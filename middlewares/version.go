package middlewares

import (
	"net/http"
	"strings"
	"user/utils/ctxlib"
	"user/utils/versionexec"

	"github.com/gin-gonic/gin"
)

type optionVersionLookup func(*gin.Context) string

type VersionOptions struct {
	VersionParamLookup optionVersionLookup
	AcceptedVersions   []string
}

// APIVersionGuard
// Middleware function to check Accept-version from API Header
func (m *Middlewares) APIVersionGuard(option VersionOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var version string

		if option.VersionParamLookup != nil {
			version = option.VersionParamLookup(ctx)
		} else {
			version = ctx.Param("version")
		}

		if version == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing version parameter"})
			return
		}

		// get and prepare the version name
		apiVersion := versionexec.PrepareVersionName(version)
		apiVersion = strings.ToUpper(apiVersion)

		var formattedVersions []string

		for _, version := range option.AcceptedVersions {
			formattedVersion := versionexec.PrepareVersionName(version)
			formattedVersions = append(formattedVersions, formattedVersion)
		}

		// set the list of system accepting version in the context
		systemAcceptedVersionsList := formattedVersions
		ctx.Set(string(ctxlib.CtxKeySystemAcceptedVersions), systemAcceptedVersionsList)

		// check the version exists in the accepted list
		// find index of version from Accepted versions
		var found bool
		for index, version := range systemAcceptedVersionsList {
			version = strings.ToUpper(version)
			if version == apiVersion {
				found = true
				ctx.Set(ctxlib.CtxKeyAcceptedVersionIndex, index)
			}
		}
		if !found {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Given version is not supported by the system"})
			return
		}

		ctx.Next()
	}
}
