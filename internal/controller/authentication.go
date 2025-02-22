package controller

import (
	"log"
	"net/http"
	"strings"
	"user/consts"
	"user/errortools"
	"user/internal/entities"
	"user/utils/apikit"
	"user/utils/stringutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (usr *User) HealthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "The server is up and running smoothly.",
	})
}

func (usr *User) Signin(ctx *gin.Context) {
	var req entities.SignInRequest
	// Bind the JSON request body to the SignInRequest struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Binding failed for user signin %v", err)
		apikit.BindingError(ctx, err)
		return
	}

	token, user, errs := usr.usecase.Signin(ctx, req)
	if !errs.Nil() {
		log.Printf("User signin failed %v", errs)
		apikit.ErrorGenerator(ctx, errs)
		return
	}

	// Return the token on successful sign-in
	response := apikit.SuccessGenerator(user, consts.User, "Signin Successful")
	response.AddData("token", token)
	ctx.JSON(http.StatusOK, response)
}

func (usr *User) Signup(ctx *gin.Context) {
	var request entities.User
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Binding failed for user signup %v", err)
		apikit.BindingError(ctx, err)
		return
	}

	request.Token = ctx.Query("token")
	// Call the use case to handle business logic (e.g., save to database)
	err := usr.usecase.SignUp(ctx, request)
	if err != nil {
		log.Printf("User signup failed %v", err)
		apikit.ErrorGenerator(ctx, err)
		return
	}

	res := apikit.SuccessGenerator(nil, "", "User signup successfully")
	ctx.JSON(http.StatusCreated, res)
}

func (usr *User) TokenValidation(ctx *gin.Context) {
	var params entities.ValidationParams

	// Bind query parameters to the struct
	if err := ctx.BindQuery(&params); err != nil {
		log.Printf("Binding failed for token validation %v", err)
		apikit.BindingError(ctx, err)
		return
	}

	// Extract token from Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		apikit.ErrorGenerator(ctx, errortools.New(errortools.UnauthorizedAccess))
		return
	}

	// Ensure the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		apikit.ErrorGenerator(ctx, errortools.New(errortools.UnauthorizedAccess))
		return
	}

	// Extract token part after "Bearer "
	params.Token = strings.TrimPrefix(authHeader, "Bearer ")

	// Print params for debugging - remember to remove or secure sensitive data in production
	permission, err := usr.usecase.ValidateToken(ctx, params)
	if err != nil {
		apikit.ErrorGenerator(ctx, err)
		return
	}

	if !params.Permission {
		res := apikit.SuccessGenerator(permission.UserID, consts.UserID, "Token validation successful")
		ctx.JSON(http.StatusOK, res)
		return
	}
	// Return the permission map
	response := apikit.SuccessGenerator(
		permission, "permissions", "Validation successful, permissions retrieved successfully")
	ctx.JSON(http.StatusOK, response)
}

func (usr *User) AddUserPermissions(ctx *gin.Context) {
	var request entities.UserProjectPermissions

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Failed to bind JSON for AddUserPermissions: %v", err)
		apikit.BindingError(ctx, err)
		return
	}

	uniqueUUIDs := stringutil.RemoveDuplicates(request.PermissionIDs)

	// Assign the unique UUIDs back to the request
	request.PermissionIDs = uniqueUUIDs

	// Add user permissions via the use case layer
	inviteLink, err := usr.usecase.AddUserPermissions(ctx, request)
	if err != nil {
		log.Printf("Error adding user permissions: %v", err)
		apikit.ErrorGenerator(ctx, err)
		return
	}

	// Prepare response based on invite link presence
	var response apikit.APIResponse
	if inviteLink != "" {
		response = apikit.SuccessGenerator(inviteLink, consts.InvitationLink, "User invited successfully")
	} else {
		response = apikit.SuccessGenerator(nil, "", "User permissions added successfully")
	}

	ctx.JSON(http.StatusOK, response)
}

func (usr *User) ListAllPermissions(ctx *gin.Context) {
	// Fetch all permissions from the use case
	permissions, errs := usr.usecase.ListAllPermissions(ctx)
	if errs != nil {
		log.Printf("Error fetching permissions: %v", errs)
		apikit.ErrorGenerator(ctx, errs)
		return
	}

	response := apikit.SuccessGenerator(permissions.PermissionsMap, "Permissions", consts.PermissionListedSuccess)
	ctx.JSON(http.StatusOK, response)
}

func (usr *User) RefreshToken(ctx *gin.Context) {
	var request entities.RefreshToken

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Failed to bind JSON for refresh token: %v", err)
		apikit.BindingError(ctx, err)
		return
	}

	token, errs := usr.usecase.RefreshToken(ctx, request.RefreshToken)
	if !errs.Nil() {
		log.Printf("Token refresh failed %v", errs)
		apikit.ErrorGenerator(ctx, errs)
		return
	}
	response := apikit.SuccessGenerator(token, "token", "Token refresh successful")
	ctx.JSON(http.StatusOK, response)
}

func (usr *User) ForgotPassword(ctx *gin.Context) {
	var request entities.ForgotPassword

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Failed to bind JSON for refresh token: %v", err)
		apikit.BindingError(ctx, err)
		return
	}

	resetLink, errs := usr.usecase.ForgotPassword(ctx, request.MailID)
	if !errs.Nil() {
		log.Printf("failed to send password reset link %v", errs)
		apikit.ErrorGenerator(ctx, errs)
		return
	}
	response := apikit.SuccessGenerator(resetLink,
		"reset_link", "successfully send reset password link to the provided mail")
	ctx.JSON(http.StatusOK, response)
}

func (usr *User) ResetPassword(ctx *gin.Context) {
	var request entities.ResetPassword

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Failed to bind JSON for reset password: %v", err)
		apikit.BindingError(ctx, err)
		return
	}
	request.Token = ctx.Query("token")

	errs := usr.usecase.ResetPassword(ctx, request)
	if !errs.Nil() {
		log.Printf("password reset failed %v", errs)
		apikit.ErrorGenerator(ctx, errs)
		return
	}
	response := apikit.SuccessGenerator(nil, "", "password reset successfully")
	ctx.JSON(http.StatusOK, response)
}
func (usr *User) ListUserPermissions(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	projectID := ctx.Param("project_id")

	if userID == "" || projectID == "" {
		apikit.ErrorGenerator(ctx, errortools.New(errortools.NoRecord))
		return
	}

	if _, err := uuid.Parse(projectID); err != nil {
		apikit.ErrorGenerator(ctx, errortools.New(errortools.NoRecord))
		return
	}

	if _, err := uuid.Parse(userID); err != nil {
		apikit.ErrorGenerator(ctx, errortools.New(errortools.NoRecord))
		return
	}

	// Fetch all permissions from the use case
	permissions, errs := usr.usecase.ListUserPermissions(ctx, userID, projectID)
	if errs != nil {
		log.Printf("Error fetching permissions: %v", errs)
		apikit.ErrorGenerator(ctx, errs)
		return
	}

	// Generate success response with the list of all permissions
	response := apikit.SuccessGenerator(permissions, "Permissions", consts.PermissionListedSuccess)
	ctx.JSON(http.StatusOK, response)
}
