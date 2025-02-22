package usecase

import (
	"context"
	"user/adapter/apihook"
	"user/adapter/emailer"
	"user/config"
	"user/errortools"
	"user/internal/entities"
	"user/internal/repo"
)

type User struct {
	repo    repo.UserImply
	cfg     *config.EnvConfig
	mailSrv emailer.Emailer
	apiSrv  apihook.HTTPAPI
}

type UserImply interface {
	AddUserPermissions(context.Context, entities.UserProjectPermissions) (string, *errortools.Error)
	Signin(context.Context, entities.SignInRequest) (entities.TokenResponse, entities.UserResponse, *errortools.Error)
	SignUp(context.Context, entities.User) *errortools.Error
	CreateUser(context.Context, entities.User) *errortools.Error
	ValidateToken(context.Context, entities.ValidationParams) (*entities.UserPermission, *errortools.Error)
	RefreshToken(context.Context, string) (*entities.TokenResponse, *errortools.Error)
	ListAllPermissions(context.Context) (*entities.Permissions, *errortools.Error)
	ListUser(ctx context.Context, pagination entities.Pagination,
		filter entities.UserFilter) ([]*entities.User, int64, error)
	ForgotPassword(context.Context, string) (string, *errortools.Error)
	ResetPassword(context.Context, entities.ResetPassword) *errortools.Error
	ListUserPermissions(context.Context, string, string) ([]entities.PermissionDetail, *errortools.Error)
}

func NewUser(repo repo.UserImply, cfg *config.EnvConfig, mailSrv emailer.Emailer, apiSrv apihook.HTTPAPI) UserImply {
	return &User{
		repo:    repo,
		cfg:     cfg,
		mailSrv: mailSrv,
		apiSrv:  apiSrv,
	}
}

func (usr *User) ListUser(ctx context.Context, pagination entities.Pagination,
	filter entities.UserFilter) ([]*entities.User, int64, error) {

	users, total, err := usr.repo.ListUser(ctx, pagination, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
