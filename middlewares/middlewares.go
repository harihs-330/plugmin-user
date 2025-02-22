package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"user/config"
	"user/consts"
	"user/errortools"
	"user/internal/entities"
	"user/internal/usecase"

	"github.com/gin-gonic/gin"
)

// Middlewares structure for storing middleware values
type Middlewares struct {
	Cfg     *config.EnvConfig
	usecase usecase.UserImply
}

// NewMiddlewares creates a new middleware object
func NewMiddlewares(cfg *config.EnvConfig, usecase usecase.UserImply) *Middlewares {
	return &Middlewares{
		Cfg:     cfg,
		usecase: usecase,
	}
}

// func (m *Middlewares) ClientDetails() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// Get the client IP using the ClientIP() method
// 		clientIP := ctx.ClientIP()

// 		// Retrieve the User-Agent header for device details
// 		deviceDetails := ctx.GetHeader("User-Agent")

// 		// Store the client IP and device details in the context
// 		ctx.Set(consts.CtxClientIP, clientIP)
// 		ctx.Set(consts.CtxDeviceDetails, deviceDetails)

// 		// Proceed to the next middleware or handler
// 		ctx.Next()
// 	}
// }

var exemptPaths = map[string]struct{}{
	"/api/:version/health":              {},
	"/api/:version/users/signin":        {},
	"/api/:version/users/signup":        {},
	"/api/:version/users/refresh-token": {},
}

func (m *Middlewares) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.GetHeader("Authorization")

		if tokenHeader != "" {
			// Split the token (Bearer <token>)
			parts := strings.Split(tokenHeader, " ")
			if len(parts) != 2 {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
				return
			}

			tokenString := parts[1]

			// Set up validation params for the user service
			params := entities.ValidationParams{
				Token:      tokenString,
				Permission: true,
			}

			userPermission, err := m.usecase.ValidateToken(ctx, params)
			if err != nil {
				switch err.Code {
				case errortools.UnauthorizedAccess:
					ctx.JSON(http.StatusUnauthorized, gin.H{
						"error": "Unauthorized Access",
					})
				case errortools.TokenExpired:
					ctx.JSON(consts.StatusInvalidToken, gin.H{
						"error": "Your token has expired. Please login again.",
					})
				default:
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": fmt.Sprintf("Error validating token: %v", err),
					})
				}
				ctx.Abort()

				return
			}

			ctx.Set(consts.AccessTokenKey, tokenString)
			ctx.Set(consts.OwnerIDKey, userPermission.UserID)
			ctx.Set(consts.PermKey, userPermission.Permissions)
			ctx.Set(consts.FullPath, ctx.FullPath())
			ctx.Set(consts.MethodKey, ctx.Request.Method)
		} else {
			// If token is missing, check if the path is exempted
			if _, exempt := exemptPaths[ctx.FullPath()]; !exempt {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"Authorization": "Missing Authorization token",
				})
				ctx.Abort()

				return
			}
		}
		ctx.Next()
	}
}
