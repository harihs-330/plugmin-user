package consts

import (
	"time"
)

// AppName stores the Application name
const (
	AppName      = "plugmin_user"
	DatabaseType = "postgres"
)

const (
	InvitationTokenDuration = 24 * time.Hour   // Invitation token expires in 24 hours
	AccessTokenDuration     = 15 * time.Minute // Access token expires in 15 minutes
	RefreshTokenDuration    = 24 * time.Hour   // Refresh token expires in 24 hours
	ResetTokenDuration      = 15 * time.Minute
	AccessToken             = "access"
	RefreshToken            = "refresh"
	InvitationToken         = "invitation"
	Permissions             = "permissions"
	ResetToken              = "password_reset"
)

var ValidSorts = []string{"id", "name", "organization", "created_on"}

const (
	InvitationMailSubject    = "Invitation"
	ResetPasswordMailSubject = "Reset Password"
)

const BcryptCost = 12 // You can adjust this value as needed

const UniqueViolationErrorCode = "23505" // PostgreSQL error code for unique violation

const MethodPost = "POST"

const (
	User                    = "user"
	UserID                  = "user_id"
	InvitationLink          = "invitation_link"
	PermissionListedSuccess = "All permissions retrieved successfully"
	ProjectID               = "project_id"
	MetaData                = "meta_data"
	Page                    = "page"
	Limit                   = "limit"
	StatusActive            = true
	StatusInActive          = false
)

const (
	StatusInvalidToken = 498
)

const (
	DefaultPage      = 1
	DefaultLimit     = 10
	DefaultSortOrder = "ASC"
	MaxLimit         = 50
	DefaultSort      = "created_on"
)

const (
	Success = "success"
	Asc     = "asc"
	Desc    = "desc"
	Sort    = "sort"
	Order   = "order"
)

const (
	UserListSuccessMsg = "Users retrieved successfully"
	CtxClientIP        = "client_ip"
	CtxDeviceDetails   = "device_details"
)

const (
	OwnerIDKey     = "user_id"
	PermKey        = "permissions"
	FullPath       = "full_path"
	MethodKey      = "method_key"
	AccessTokenKey = "access_token"
)
