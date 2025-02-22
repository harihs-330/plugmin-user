package crypto

import (
	"fmt"
	"log"
	"net/url"
	"time"
	"user/internal/entities"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// InvitationTokenClaims represents the JWT claims for invitation tokens
type InvitationTokenClaims struct {
	ID            string      `json:"id"`
	Email         string      `json:"email"` // Added email field
	ProjectID     string      `json:"project_id"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
	TokenType     string      `json:"token_type"` // Token type (invitation)
	Exp           time.Time   `json:"exp"`        // Expiration time
	Iss           string      `json:"iss"`        // Issuer
	jwt.StandardClaims
}

// TokenClaims represents the JWT claims
type TokenClaims struct {
	ID     string    `json:"id"`
	UserID string    `json:"user_id"`
	Email  string    `json:"email"`
	Exp    time.Time `json:"exp"`
	Iat    time.Time `json:"iat"`
	Iss    string    `json:"iss"`
	Sub    string    `json:"sub"`
	jwt.StandardClaims
}

// NewInvitationTokenClaims creates a new InvitationTokenClaims instance with the provided values
func NewInvitationTokenClaims(
	projectID, tokenType, issuer, email string, permissionIDs []uuid.UUID, exp time.Time) *InvitationTokenClaims {

	return &InvitationTokenClaims{
		ProjectID:     projectID,
		PermissionIDs: permissionIDs,
		TokenType:     tokenType,
		Iss:           issuer,
		Exp:           exp,
		Email:         email,
	}
}

// NewTokenClaims creates a new TokenClaims instance with the provided values
func NewTokenClaims(userID, email, issuer, subject string, exp time.Time) *TokenClaims {
	return &TokenClaims{
		UserID: userID,
		Email:  email,
		Exp:    exp,
		Iss:    issuer,
		Sub:    subject,
	}
}

// CreateInvitationToken generates a JWT token for invitation purposes
func (itk *InvitationTokenClaims) CreateInvitationToken(secretKey string) (string, error) {
	// Generate a new UUID for the token ID
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// Set the unique token ID and issued at time
	itk.ID = uid.String()
	itk.StandardClaims.IssuedAt = time.Now().Unix() // Unix timestamp for issued time
	itk.StandardClaims.ExpiresAt = itk.Exp.Unix()   // Set expiration time as Unix timestamp

	// Create JWT with the claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, itk)

	// Sign the token using the secret key
	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("[CreateInvitationToken] Error signing token: %v", err)
		return "", err
	}

	// Return the created invitation token
	return tokenString, nil
}

// GenerateInvitationLink generates an invitation URL with the token as a query parameter
func GenerateInvitationLink(baseURL, token string) (string, error) {
	// Parse the base URL to ensure it's a valid URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Append the token as a query parameter
	queryParams := parsedURL.Query()
	queryParams.Set("token", token)
	parsedURL.RawQuery = queryParams.Encode()

	// Return the full invitation URL
	return parsedURL.String(), nil
}

// CreateToken generates a JWT token and returns a Token struct with the token details
func (tk *TokenClaims) CreateToken(secretKey string) (entities.Token, error) {
	// Generate a new UUID for the token ID
	uid, err := uuid.NewRandom()
	if err != nil {
		return entities.Token{}, err
	}

	// Set the token ID in the claims
	tk.ID = uid.String()
	tk.Iat = time.Now().UTC()

	// Create the JWT claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

	// Sign the token with the provided secret key
	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("[CreateToken] Error signing token: %v", err)
		return entities.Token{}, err
	}

	// Return the token struct with generated values
	return entities.Token{
		ID:        uid.String(),
		UserID:    tk.UserID,
		Token:     tokenString,
		TokenType: tk.Sub, // TokenType is used to store the token type
		CreatedOn: time.Now().UTC(),
		ExpiresOn: tk.Exp,
	}, nil
}

// HashPassword hashes the given password using bcrypt with a specified cost
func HashPassword(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

// CheckPassword compares the hashed password with the plain text password
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func ExtractTokenClaims(tokenString, secretKey string, claims jwt.Claims) (jwt.Claims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("[ExtractTokenClaims] Error parsing token: %v", err)
		return nil, err
	}

	// Validate the token and ensure claims are valid
	if token.Valid {
		return token.Claims, nil
	}

	log.Println("[ExtractTokenClaims] Invalid token or claims")

	return nil, fmt.Errorf("invalid token")
}
