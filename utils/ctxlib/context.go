package ctxlib

import "context"

// use this type to pass around contexts
type CtxKey string

const (
	CtxKeySystemAcceptedVersions CtxKey = "System-Accept-Versions"
	CtxKeyAcceptedVersionIndex          = "Accepted-Version-Index"
	CtxKeyOwnerID                CtxKey = "OwnerID"
)

// get Context
// read the context data and type assert into corresponding concrete value
func Get[T any](ctx context.Context, name string) (T, bool) {
	value, exists := ctx.Value(name).(T)

	return value, exists
}
