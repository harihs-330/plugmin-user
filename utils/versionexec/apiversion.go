package versionexec

import (
	"fmt"
	"reflect"
	"strings"
	"user/utils/ctxlib"

	"github.com/gin-gonic/gin"
)

// RenderHandler
// the handler method should always check version_method exists or not
// if that exists, it will execute it, instead the given method
func RenderHandler(ctx *gin.Context, object interface{}, method string, args ...interface{}) {
	inputs := make([]reflect.Value, 0, len(args))

	// passing the context to the methods
	// first argument should be the ctx
	inputs = append(inputs, reflect.ValueOf(ctx))

	for _, v := range args {
		inputs = append(inputs, reflect.ValueOf(v))
	}

	// read the data from context
	systemAcceptedVs, _ := ctxlib.Get[[]string](ctx, string(ctxlib.CtxKeySystemAcceptedVersions))
	headerVersionIndex, _ := ctxlib.Get[int](ctx, ctxlib.CtxKeyAcceptedVersionIndex)

	// loop thorugh
	for i := len(systemAcceptedVs[0:headerVersionIndex]); i >= 0; i-- {
		versionMethod := fmt.Sprintf("%s_%s", strings.ToUpper(systemAcceptedVs[i]), method)

		// check object implement the method
		// like if the method is GetUsers, and version is v1 ; it will check v1_GetUsers
		callableMethod := reflect.ValueOf(object).MethodByName(versionMethod)
		if callableMethod.IsValid() {
			// callableMethod.Call(inputs)[0].Interface()
			callableMethod.Call(inputs)
			return
		}
	}

	// check objConv implement the method
	callableMethod := reflect.ValueOf(object).MethodByName(method)
	if callableMethod.IsValid() {
		// callableMethod.Call(inputs)[0].Interface()
		callableMethod.Call(inputs)

		return
	}

	panic(fmt.Sprintf("unable to locate the method %v", method))
}

// PrepareVersionName
// this function will prepare a version name v1.1 into v1_1
func PrepareVersionName(version string) string {
	version = strings.ReplaceAll(version, ".", "_")
	return version
}
