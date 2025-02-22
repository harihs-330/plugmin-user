package entities

import (
	"strings"
	"time"
	"user/consts"
	"user/errortools"

	"github.com/google/uuid"
)

type Database struct {
	DBName   string `json:"db_name"`
	Driver   string `json:"driver"`
	User     string `json:"user_name"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Schema   string `json:"schema"`
	SSLMode  string `json:"ssl_mode"`
}

type UserProjectPermissions struct {
	MailID        string      `json:"email" validate:"required,email"`
	ProjectID     string      `json:"project_id" validate:"required,trimmedempty,uuid"`
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"required,trimmedempty"`
	UserID        string
}

// User represents a user entity with fields for user management
type User struct {
	UserID          string
	Name            string     `json:"name" validate:"required,trimmedempty,min=2,max=100,customname"`
	Password        string     `json:"password,omitempty" validate:"required,passwordformat"`
	ConfirmPassword string     `json:"confirm_password,omitempty" validate:"required,eqfield=Password"`
	Email           string     `json:"email" validate:"required,customemail" `
	Purpose         string     `json:"purpose" validate:"required,trimmedempty,oneof=student developer other"`
	Organization    string     `json:"organization" validate:"required,trimmedempty,min=2,max=100,customname"`
	HashedPassword  string     `json:"hashed_password,omitempty"`
	CreatedOn       time.Time  `json:"created_on,omitempty"`
	UpdatedOn       time.Time  `json:"updated_on,omitempty"`
	IsDeleted       bool       `json:"is_deleted,omitempty"`
	ISActive        bool       `json:"is_active"`
	DeletedOn       *time.Time `json:"deleted_on,omitempty"`
	Token           string     `json:"token,omitempty"`
}

type UserFilter struct {
	ID           string `form:"id"`
	Name         string `form:"name"`
	Email        string `form:"mailid"`
	Organization string `form:"organization"`
	ProjectID    string `form:"project_id"`
	IsActive     string `form:"is_active"`
	IsMember     string `form:"is_member"`
	Pagination   Pagination
}

var OneOfValuesMap = map[string][]string{
	"purpose": {"student", "developer", "other"},
}

type SignInRequest struct {
	MailID   string ` json:"mail_id" validate:"required,customemail"`
	Password string `binding:"required" json:"password"`
}

type UserResponse struct {
	UserID       string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Purpose      string `json:"purpose"`
	Organization string `json:"organization"`
}

type Token struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	CreatedOn time.Time `json:"created_on"`
	ExpiresOn time.Time `json:"expires_on"`
}

type TokenResponse struct {
	AccessToken       string    `json:"access_token"`
	RefreshToken      string    `json:"refresh_token"`
	AccessTokenExpiry time.Time `json:"access_token_expiry"`
}

type ValidationParams struct {
	Permission bool   `form:"permission"`
	ProjectID  string `form:"project_id"`
	Token      string
}

type UserPermission struct {
	UserID      string                       `json:"user_id"`
	Permissions map[string]map[string]string `json:"permissions"`
}

type PermissionResp struct {
	UserID      string                       `json:"user_id"`
	Permissions map[string]map[string]string `json:"permissions"`
}

type PermissionDetail struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Permissions struct {
	PermissionsMap map[uuid.UUID]PermissionDetail `json:"permissions"`
}
type Pagination struct {
	Page  int    `form:"page" `
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
	Order string `form:"order"`
}

type MetaData struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	Next        int `json:"next"`
	Prev        int `json:"prev"`
}

//nolint:cyclop
func (pg *Pagination) Validate(validSorts ...string) *errortools.Error {
	validationErr := errortools.Init()
	if pg.Page == 0 {
		pg.Page = consts.DefaultPage
	}
	if pg.Limit == 0 {
		pg.Limit = consts.DefaultLimit
	}
	if pg.Page <= 0 {
		validationErr.AddValidationError(consts.Page, errortools.InvalidPageErrorCode)
	}

	if pg.Limit <= 0 || pg.Limit > consts.MaxLimit {
		validationErr.AddValidationError(consts.Limit, errortools.InvalidLimitErrorCode, consts.MaxLimit)
	}

	if pg.Sort != "" {
		valid := false
		for _, v := range validSorts {
			if v == pg.Sort {
				valid = true
				break
			}
		}
		if !valid {
			validSortKeysMessage := strings.Join(validSorts, ", ")
			validationErr.AddValidationError(consts.Sort, errortools.InvalidSortBy, validSortKeysMessage)
		}
	}

	if pg.Order != "" {
		orderLower := strings.ToLower(pg.Order)
		if orderLower != consts.Asc && orderLower != consts.Desc {
			validationErr.AddValidationError(consts.Order, errortools.InvalidSortOrder, pg.Order)
		}
	}

	return validationErr
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPassword struct {
	MailID string `json:"mail_id"`
}

type ResetPassword struct {
	Password        string `json:"password,omitempty" validate:"required,passwordformat"`
	ConfirmPassword string `json:"confirm_password,omitempty" validate:"required,eqfield=Password"`
	Token           string `json:"token,omitempty"`
	HashedPassword  string `json:"hashed_password,omitempty"`
	UserID          string
}
