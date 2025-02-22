package controller

import (
	"log"
	"net/http"
	"user/config"
	"user/consts"
	"user/internal/entities"
	"user/internal/usecase"
	"user/utils"
	"user/utils/apikit"
	"user/utils/stringutil"
	"user/utils/versionexec"

	"user/errortools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	router  *gin.RouterGroup
	usecase usecase.UserImply
	cfg     *config.EnvConfig
}

// NewUser creates a new User instance.
func NewUser(router *gin.RouterGroup, usecase usecase.UserImply, cfg *config.EnvConfig) *User {
	return &User{
		router:  router,
		usecase: usecase,
		cfg:     cfg,
	}
}

// InitRoutes initializes the routes for the User controller.
func (usr *User) InitRoutes() {
	usr.router.GET("/health", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "HealthHandler")
	})

	usr.router.POST("/signin", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "Signin")
	})
	usr.router.POST("/signup", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "Signup")
	})
	usr.router.GET("/token/validate", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "TokenValidation")
	})
	usr.router.PUT("/permissions", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "AddUserPermissions")
	})
	usr.router.GET("/permissions", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "ListAllPermissions")
	})
	usr.router.POST("/refresh-token", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "RefreshToken")
	})
	usr.router.GET("/", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "ListUsers")
	})
	usr.router.POST("/forgot-password", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "ForgotPassword")
	})
	usr.router.POST("/reset-password", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "ResetPassword")
	})
	usr.router.GET("/projects/:project_id/users/:user_id/permissions", func(ctx *gin.Context) {
		versionexec.RenderHandler(ctx, usr, "ListUserPermissions")
	})
}

func (usr *User) ListUsers(ctx *gin.Context) {
	var (
		filter        entities.UserFilter
		validationErr = errortools.Init()
	)

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		apikit.ErrorGenerator(ctx, err)
		return
	}

	if filter.ProjectID != "" {
		if _, err := uuid.Parse(filter.ProjectID); err != nil {
			validationErr.AddValidationError(
				consts.ProjectID,
				errortools.InvalidUUID,
				consts.ProjectID,
			)
			apikit.ErrorGenerator(ctx, validationErr)

			return
		}
	}

	filter.Pagination.Page = stringutil.ParseQueryParam(ctx.Query("page"), consts.DefaultPage).ToInt()
	filter.Pagination.Limit = stringutil.ParseQueryParam(ctx.Query("limit"), consts.DefaultLimit).ToInt()
	filter.Pagination.Sort = stringutil.ParseQueryParam(ctx.Query("sort"), consts.DefaultSort).ToString()
	filter.Pagination.Order = stringutil.ParseQueryParam(ctx.Query("order"), consts.DefaultSortOrder).ToString()

	if validationErr = filter.Pagination.Validate(consts.ValidSorts...); !validationErr.Nil() {
		apikit.ErrorGenerator(ctx, validationErr)
		return
	}

	users, total, err := usr.usecase.ListUser(ctx, filter.Pagination, filter)
	if err != nil {
		log.Printf("failed to get users: %v", err)
		apikit.ErrorGenerator(ctx, err)
		return
	}

	if len(users) == 0 {
		users = make([]*entities.User, 0)
	}

	metaData := utils.MetaDataInfo(&entities.MetaData{
		Total:       int(total),
		PerPage:     filter.Pagination.Limit,
		CurrentPage: filter.Pagination.Page,
	})

	response := apikit.SuccessGenerator(users, consts.User, consts.UserListSuccessMsg)
	if total > 0 {
		response.AddData(consts.MetaData, metaData)
	}
	ctx.JSON(http.StatusOK, response)
}
