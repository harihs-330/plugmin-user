package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"user/adapter/apihook"
	"user/consts"
	"user/errortools"
	"user/internal/entities"
	"user/utils/crypto"
	"user/utils/ctxlib"
	"user/utils/parser"
	"user/utils/sanitizer"

	"user/utils/structutil"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func (usr *User) Signin(ctx context.Context,
	request entities.SignInRequest,
) (entities.TokenResponse, entities.UserResponse, *errortools.Error) {

	validationErr := errortools.Init()
	validate := sanitizer.Validator()
	// Validate the schema structure
	if err := validate.Struct(&request); err != nil {
		return entities.TokenResponse{}, entities.UserResponse{}, processValidationErrors(err, validationErr)
	}
	userData, err := usr.repo.FetchUserByEmail(ctx, request.MailID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("No user exists for the given mailID %s", request.MailID)
		return entities.TokenResponse{}, entities.UserResponse{}, errortools.New(errortools.InvalidEmailID)
	}
	if err != nil {
		log.Printf("Signin: failed to fetch user by email %s: %v", request.MailID, err)
		return entities.TokenResponse{}, entities.UserResponse{}, errortools.New(errortools.InternalServerErrCode)
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(userData.HashedPassword), []byte(request.Password)); err != nil {
		log.Printf("Signin: password mismatch for user %s: %v", userData.Email, err)
		return entities.TokenResponse{}, entities.UserResponse{}, errortools.New(errortools.WrongPassword)
	}

	// Create access token claims
	accessTokenClaims := crypto.NewTokenClaims(
		userData.UserID,
		userData.Email,
		consts.AppName,
		consts.AccessToken,
		time.Now().UTC().Add(consts.AccessTokenDuration),
	)

	// Create access token
	accessToken, err := accessTokenClaims.CreateToken(usr.cfg.SecretKey)
	if err != nil {
		log.Printf("Signin: failed to create access token for user %s: %v", userData.Email, err)
		return entities.TokenResponse{}, entities.UserResponse{}, errortools.New(errortools.InternalServerErrCode)
	}

	// Create refresh token claims
	refreshTokenClaims := crypto.NewTokenClaims(
		userData.UserID,
		userData.Email,
		consts.AppName,
		consts.RefreshToken,
		time.Now().UTC().Add(consts.RefreshTokenDuration),
	)

	// Create refresh token
	refreshToken, err := refreshTokenClaims.CreateToken(usr.cfg.SecretKey)
	if err != nil {
		log.Printf("Signin: failed to create refresh token for user %s: %v", userData.Email, err)
		return entities.TokenResponse{}, entities.UserResponse{}, errortools.New(errortools.InternalServerErrCode)
	}

	// Insert tokens into the database
	tokens := []entities.Token{accessToken, refreshToken}
	columnCount := structutil.GetStructFieldCount(refreshToken)
	if err := usr.repo.InsertTokens(ctx, tokens, columnCount); err != nil {
		log.Printf("Signin: failed to insert tokens for user %s: %v", userData.Email, err)
		return entities.TokenResponse{}, entities.UserResponse{}, errortools.New(errortools.InternalServerErrCode)
	}

	// Return the response with tokens and user data
	return entities.TokenResponse{
			AccessToken:       accessToken.Token,
			RefreshToken:      refreshToken.Token,
			AccessTokenExpiry: accessToken.ExpiresOn,
		}, entities.UserResponse{
			UserID:       userData.UserID,
			Name:         userData.Name,
			Email:        userData.Email,
			Purpose:      userData.Purpose,
			Organization: userData.Organization,
		}, nil
}

func (usr *User) SignUp(ctx context.Context, request entities.User) *errortools.Error { //nolint
	// Check if the token is provided for invitation-based registration
	if request.Token == "" {
		// If no token, directly create the user
		if err := usr.CreateUser(ctx, request); err != nil {
			return err
		}
		return nil
	}

	// Handle invitation token logic
	invitationClaims := &crypto.InvitationTokenClaims{}
	claims, err := crypto.ExtractTokenClaims(request.Token, usr.cfg.SecretKey, invitationClaims)
	if err != nil {
		log.Printf("Failed to extract details from token: %v", err)
		return errortools.New(errortools.UnauthorizedAccess)
	}
	extractedDetails, ok := claims.(*crypto.InvitationTokenClaims)
	if !ok {
		log.Printf("Failed to convert claims: %v", err)
		return errortools.New(errortools.UnauthorizedAccess)
	}

	// Validate email in the token and request
	if extractedDetails.Email != request.Email {
		log.Printf("Token email: %s and request email: %s do not match", extractedDetails.Email, request.Email)
		return errortools.New(errortools.UnauthorizedAccess)
	}

	// Check if the invitation token has expired
	if time.Now().After(extractedDetails.Exp) {
		log.Printf("Invitation token for user has expired %s", request.Email)
		return errortools.New(errortools.TokenExpired)
	}

	// Create user
	if err := usr.CreateUser(ctx, request); err != nil {
		return err
	}

	log.Printf("User created successfully: %s", request.Email)

	requestBody := entities.UserProjectPermissions{
		MailID:        request.Email,
		PermissionIDs: extractedDetails.PermissionIDs,
		ProjectID:     extractedDetails.ProjectID,
	}
	_, errs := usr.AddUserPermissions(ctx, requestBody)
	if errs != nil {
		log.Printf("Failed to add user permissions: %v", errs)
		return errs
	}

	newToken, resp, fnErr := usr.Signin(ctx, entities.SignInRequest{
		MailID:   request.Email,
		Password: request.Password,
	})

	if fnErr != nil {
		log.Println("Error while calling signin", fnErr)
		return &errortools.Error{Msg: fnErr.Error()}
	}

	apiErr := usr.UserProjectMapping(ctx, extractedDetails.ProjectID, resp.UserID, newToken.AccessToken)
	if apiErr != nil {
		return apiErr
	}
	// Send the welcome email after successful user creation
	return nil
}

func (usr *User) sendWelcomeEmail(email, name string) *errortools.Error {
	emailData := map[string]string{
		"UserEmail": email,
		"UserName":  name,
	}

	// Render email template
	emailContent, err := parser.HTML("templates/signup_success.html", emailData)
	if err != nil {
		log.Printf("Failed to render email template: %v", err)
		return errortools.New(errortools.InternalServerErrCode)
	}

	// Log the rendered content to ensure the template is correct
	log.Printf("Rendered Email Content: %s", emailContent)

	// Send email
	err = usr.mailSrv.Send(
		"Welcome to Plugmin",
		emailContent,
		[]string{email}, // To list
		nil,             // CC list
		nil,             // BCC list
		nil,             // Attachments
	)
	if err != nil {
		log.Printf("Failed to send welcome email: %v", err)
		return errortools.New(errortools.InternalServerErrCode)
	}

	log.Printf("Welcome email sent successfully to %s", email)

	return nil
}

func (usr *User) CreateUser(ctx context.Context, request entities.User) *errortools.Error {
	validationErr := errortools.Init()
	validate := sanitizer.Validator()

	// Validate the schema structure
	if err := validate.Struct(&request); err != nil {
		return processValidationErrors(err, validationErr)
	}

	// Hash the password using bcrypt
	hashedPassword, err := crypto.HashPassword(request.Password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return errortools.New(errortools.InternalServerErrCode)
	}

	// Set the hashed password in the request object
	request.HashedPassword = hashedPassword

	// Insert the user data into the database
	if err := usr.repo.CreateUser(ctx, request); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// Handle duplicate key violation
			if pqErr.Code == consts.UniqueViolationErrorCode {
				log.Printf("Duplicate entry error: %v", err)
				return errortools.New(errortools.DuplicateRecord)
			}
		}
		log.Printf("Failed to create user in database: %v", err)

		return errortools.New(errortools.InternalServerErrCode)
	}

	return usr.sendWelcomeEmail(request.Email, request.Name)
}

func processValidationErrors(err error, validationErr *errortools.Error) *errortools.Error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, valErr := range validationErrors {
			customMessage, exists := errortools.CustomMessages[valErr.Tag()]
			if valErr.Tag() == "oneof" {
				allowedValues := entities.OneOfValuesMap[valErr.Field()]
				valuesStr := strings.Join(allowedValues, ", ")
				customMessage = fmt.Sprintf(customMessage, valErr.Field(), valuesStr)
			}
			if strings.Contains(customMessage, "%s") {
				args := []interface{}{valErr.Field(), valErr.Param()}
				customMessage = formatCustomMessage(customMessage, args...)
			}
			if !exists {
				customMessage = fmt.Sprintf("invalid field value for %s", valErr.Field())
			}
			validationErr.AddValidationError(valErr.Field(), errortools.ValidatorErrors, customMessage)
		}
	}
	if !validationErr.Nil() {
		log.Printf("Validation errors: %v", validationErr)
		return validationErr
	}

	return nil
}
func formatCustomMessage(customMessage string, args ...interface{}) string {
	placeholderCount := strings.Count(customMessage, "%s")

	if len(args) > placeholderCount {
		args = args[:placeholderCount]
	}

	return fmt.Sprintf(customMessage, args...)
}

func (usr *User) ValidateToken(
	ctx context.Context, params entities.ValidationParams) (*entities.UserPermission, *errortools.Error) {

	token := entities.Token{
		Token:     params.Token,
		TokenType: consts.AccessToken,
	}

	userID, expiresOn, err := usr.repo.ValidateToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Token not found in the database, log and return unauthorized error
			log.Printf("Token not found: %s", params.Token)
			return nil, errortools.New(errortools.UnauthorizedAccess)
		}
		// Log and return internal server error for other issues
		log.Printf("Error during token validation: %v", err)

		return nil, errortools.New(errortools.InternalServerErrCode)
	}
	if time.Now().After(expiresOn) {
		log.Printf("Token for userID %s has expired at %v", userID, expiresOn)
		return nil, errortools.New(errortools.TokenExpired)
	}
	userperm := entities.UserPermission{}
	userperm.UserID = userID
	// Step 2: If no permission check is required, return validation success
	if !params.Permission {
		// Return an empty map and no error to indicate successful validation without permission retrieval
		return &userperm, nil
	}

	// Step 3: Get permissions if the permission flag is true
	permissions, err := usr.repo.GetPermissions(ctx, userID, params.ProjectID)
	if err != nil {
		// Log and return internal server error on permission retrieval failure
		log.Printf("Error during permissions retrieval for userID %s: %v", userID, err)
		return &userperm, errortools.New(errortools.InternalServerErrCode)
	}
	userperm.Permissions = permissions
	// Return the permissions
	return &userperm, nil
}

func (usr *User) AddUserPermissions(
	ctx context.Context, request entities.UserProjectPermissions,
) (string, *errortools.Error) {

	validationErr := errortools.Init()

	validate := sanitizer.Validator()

	// Validate the request struct
	if err := validate.Struct(&request); err != nil {
		return "", processValidationErrors(err, validationErr)
	}

	// Fetch user details by email
	userData, err := usr.repo.FetchUserByEmail(ctx, request.MailID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("User not found: %s", request.MailID)

		// Create and send invitation if the user is not found
		return usr.handleInvitation(request)
	}

	if err != nil {
		// Log other errors during user fetch
		log.Printf("Error fetching user details for email %s: %v", request.MailID, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	// Set the user ID and add permissions
	request.UserID = userData.UserID
	if err := usr.repo.AddUserPermissions(ctx, request); err != nil {
		log.Printf("Failed to add permissions for Mail ID %s: %v", request.MailID, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	token, _ := ctxlib.Get[string](ctx, consts.AccessTokenKey)
	apiErr := usr.UserProjectMapping(ctx, request.ProjectID, request.UserID, token)
	if apiErr != nil {
		return "", apiErr
	}

	return "", nil
}

func (usr *User) handleInvitation(request entities.UserProjectPermissions) (string, *errortools.Error) {
	// Create the invitation token claims
	invitationTokenClaims := crypto.NewInvitationTokenClaims(
		request.ProjectID,
		consts.InvitationToken,
		consts.AppName,
		request.MailID,
		request.PermissionIDs,
		time.Now().UTC().Add(consts.InvitationTokenDuration),
	)

	// Create the invitation token
	invitationToken, err := invitationTokenClaims.CreateInvitationToken(usr.cfg.SecretKey)
	if err != nil {
		log.Printf("Failed to create invitation token for email %s: %v", request.MailID, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	// Generate the invitation link
	invitationLink, err := crypto.GenerateInvitationLink(usr.cfg.SignupURL, invitationToken)
	if err != nil {
		log.Printf("Failed to generate invitation link for email %s: %v", request.MailID, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	// Prepare the email data
	emailData := map[string]string{
		"SignupLink": invitationLink,
	}

	// Generate the HTML content using the template
	emailContent, err := parser.HTML("templates/invite_user.html", emailData)
	if err != nil {
		log.Printf("Failed to render email HTML template: %v", err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	// Send the invitation email using the Send method
	err = usr.mailSrv.Send(
		consts.InvitationMailSubject,
		emailContent,
		[]string{request.MailID}, // To list
		[]string{},               // CC list (empty slice)
		[]string{},               // BCC list (empty slice)
		[]string{},               // Attachments list (empty slice)
	)
	if err != nil {
		log.Printf("Failed to send invitation mail for %s: %v", request.MailID, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	return invitationLink, nil
}

func (usr *User) ListAllPermissions(ctx context.Context) (*entities.Permissions, *errortools.Error) {
	permissionList, err := usr.repo.ListAllPermissions(ctx)
	if err != nil {
		log.Printf("Failed to list permissions :%v", err)
		return nil, errortools.New(errortools.FailedGettingPermission)
	}

	return permissionList, nil
}

func (usr *User) RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenResponse, *errortools.Error) {
	token := entities.Token{
		Token:     refreshToken,
		TokenType: consts.RefreshToken,
	}
	userID, expire, err := usr.repo.ValidateToken(ctx, token)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("No user exists for the given token %s", refreshToken)
		return nil, errortools.New(errortools.UnauthorizedAccess)
	}
	if err != nil {
		return nil, errortools.New(errortools.InternalServerErrCode)
	}
	if time.Now().After(expire) {
		log.Printf("Refresh Token for userID %s has expired at %v", userID, expire)
		return nil, errortools.New(errortools.TokenExpired)
	}

	invitationClaims := &crypto.TokenClaims{}
	claims, err := crypto.ExtractTokenClaims(refreshToken, usr.cfg.SecretKey, invitationClaims)
	if err != nil {
		log.Printf("Failed to extract details form token: %v", err)
		return nil, errortools.New(errortools.UnauthorizedAccess)
	}
	extractedDetails, ok := claims.(*crypto.TokenClaims)
	if !ok {
		log.Printf("Failed to convert claims: %v", err)
		return nil, errortools.New(errortools.UnauthorizedAccess)
	}
	accessTokenClaims := crypto.NewTokenClaims(
		userID,
		extractedDetails.Email,
		consts.AppName,
		consts.AccessToken,
		time.Now().UTC().Add(consts.AccessTokenDuration),
	)

	// Create access token
	accessToken, err := accessTokenClaims.CreateToken(usr.cfg.SecretKey)
	if err != nil {
		log.Printf("RefreshToken: failed to create access token for user %s: %v", extractedDetails.Email, err)
		return nil, errortools.New(errortools.InternalServerErrCode)
	}

	columnCount := structutil.GetStructFieldCount(accessToken)
	err = usr.repo.InsertTokens(ctx, []entities.Token{accessToken}, columnCount)
	if err != nil {
		log.Printf("Signin: failed to insert token for user %s: %v", extractedDetails.Email, err)
		return nil, errortools.New(errortools.InternalServerErrCode)
	}

	return &entities.TokenResponse{
		AccessToken:       accessToken.Token,
		RefreshToken:      refreshToken,
		AccessTokenExpiry: accessToken.ExpiresOn,
	}, nil
}

func (usr *User) ForgotPassword(ctx context.Context, mailID string) (string, *errortools.Error) {
	user, err := usr.repo.FetchUserByEmail(ctx, mailID)
	if errors.Is(err, sql.ErrNoRows) {
		// Token not found in the database, log and return unauthorized error
		log.Printf("User not found: %s", mailID)
		return "", errortools.New(errortools.NoRecord)
	}

	// Create access token claims
	passwordResetTokenClaims := crypto.NewTokenClaims(
		user.UserID,
		user.Email,
		consts.AppName,
		consts.ResetToken,
		time.Now().UTC().Add(consts.ResetTokenDuration),
	)

	// Create access token
	passwordResetToken, err := passwordResetTokenClaims.CreateToken(usr.cfg.SecretKey)
	if err != nil {
		log.Printf("Signin: failed to create password reset token for user %s: %v", user.Email, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	columnCount := structutil.GetStructFieldCount(passwordResetToken)
	err = usr.repo.InsertTokens(ctx, []entities.Token{passwordResetToken}, columnCount)
	if err != nil {
		log.Printf("Signin: failed to insert token for user %s: %v", user.Email, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	// Generate reset link
	resetLink, err := crypto.GenerateInvitationLink(usr.cfg.ResetPasswordURL, passwordResetToken.Token)
	if err != nil {
		log.Printf("Failed to generate password reset link for email %s: %v", user.Email, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}
	// Prepare the email data
	emailData := map[string]string{
		"ResetLink": resetLink,
	}

	// Generate the HTML content using the template
	emailContent, err := parser.HTML("templates/reset_password.html", emailData)
	if err != nil {
		log.Printf("Failed to render email HTML template: %v", err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	// Send the invitation email using the Send method
	err = usr.mailSrv.Send(
		consts.ResetPasswordMailSubject,
		emailContent,
		[]string{user.Email}, // To list
		[]string{},           // CC list (empty slice)
		[]string{},           // BCC list (empty slice)
		[]string{},           // Attachments list (empty slice)
	)

	if err != nil {
		log.Printf("Failed to send password reset mail for %s: %v", user.Email, err)
		return "", errortools.New(errortools.InternalServerErrCode)
	}

	log.Printf("Successfully sent password reset email to: %s", user.Email)

	return resetLink, nil
}

func (usr *User) ResetPassword(ctx context.Context, request entities.ResetPassword) *errortools.Error { //nolint
	validationErr := errortools.Init()
	validate := sanitizer.Validator()

	// Validate the schema structure
	if err := validate.Struct(&request); err != nil {
		return processValidationErrors(err, validationErr)
	}
	if request.Token == "" {
		return errortools.New(errortools.InvalidToken)
	}

	resetTokenClaims := &crypto.TokenClaims{}
	claims, err := crypto.ExtractTokenClaims(request.Token, usr.cfg.SecretKey, resetTokenClaims)
	if err != nil {
		log.Printf("Failed to extract details form token: %v", err)
		return errortools.New(errortools.InvalidToken)
	}
	extractedDetails, ok := claims.(*crypto.TokenClaims)
	if !ok {
		log.Printf("Failed to convert claims: %v", err)
		return errortools.New(errortools.InvalidToken)
	}
	if time.Now().After(extractedDetails.Exp) {
		log.Printf("Token for mailID %s has expired at %v", extractedDetails.Email, extractedDetails.Exp)
		return errortools.New(errortools.TokenExpired)
	}
	isRevoked, err := usr.repo.IsTokenRevoked(ctx, extractedDetails.ID)
	if isRevoked {
		log.Printf("Token for mailID %s has been revoked", extractedDetails.Email)
		return errortools.New(errortools.TokenExpired)
	}
	if errors.Is(err, sql.ErrNoRows) {
		// Token not found in the database, log and return unauthorized error
		log.Printf("Token for mailID %s not found in database", extractedDetails.Email)
		return errortools.New(errortools.InvalidToken)
	}
	if err != nil {
		log.Printf("Failed to fetch token details for mailID: %s %v", extractedDetails.Email, err)
		return errortools.New(errortools.InternalServerErrCode)
	}
	userData, err := usr.repo.FetchUserByEmail(ctx, extractedDetails.Email)
	if errors.Is(err, sql.ErrNoRows) {
		// Token not found in the database, log and return unauthorized error
		log.Printf("User not found: %s", extractedDetails.Email)
		return errortools.New(errortools.UnauthorizedAccess)
	}
	if err != nil {
		log.Printf("Failed to fetch user details for mailID: %s %v", extractedDetails.Email, err)
		return errortools.New(errortools.InternalServerErrCode)
	}
	hashedPassword, err := crypto.HashPassword(request.Password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return errortools.New(errortools.InternalServerErrCode)
	}
	request.HashedPassword = hashedPassword
	request.UserID = userData.UserID
	err = usr.repo.ResetPassword(ctx, request)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		return errortools.New(errortools.InternalServerErrCode)
	}
	err = usr.repo.RevokeToken(ctx, extractedDetails.ID)
	if err != nil {
		log.Printf("Failed to revoke token: %v", err)
		return errortools.New(errortools.InternalServerErrCode)
	}

	return nil
}

func (usr *User) UserProjectMapping(ctx context.Context,
	projectID, userID, token string) *errortools.Error {

	apiParams := apihook.Params{
		URL: usr.cfg.APIBuilderURL + "/projects/invited-users",
		Header: map[string]interface{}{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		Body: map[string]interface{}{
			"project_id": projectID,
			"user_id":    userID,
			"created_at": time.Now(),
		},
	}

	var apiResponse interface{}
	statusCode, err := usr.apiSrv.PUT(ctx, apiParams, &apiResponse)

	if err != nil {
		if statusCode != http.StatusOK {
			log.Printf("Failed to call apibuilder service: %v, Status Code: %d", err, statusCode)
			return &errortools.Error{Msg: err.Error()}
		}
		log.Println("Error while calling apibuilder", err)

		return &errortools.Error{Msg: err.Error()}
	}

	return nil
}

func (usr *User) ListUserPermissions(
	ctx context.Context,
	userID string,
	projectID string,
) ([]entities.PermissionDetail, *errortools.Error) {

	permissionList, err := usr.repo.ListUserPermissions(ctx, userID, projectID)
	if err != nil {
		log.Printf("Failed to list permissions :%v", err)
		return nil, errortools.New(errortools.FailedGettingPermission)
	}
	if len(permissionList) == 0 {
		return []entities.PermissionDetail{}, nil
	}

	return permissionList, nil
}
