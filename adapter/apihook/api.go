package apihook

import (
	"context"
	"errors"
)

var errUnimplemented = errors.New("not yet implemented")

type Params struct {
	URL          string
	Method       string
	Body, Header map[string]interface{}
}

type Connector interface {
	GET(ctx context.Context, params Params, result interface{}) (int, error)
	POST(ctx context.Context, params Params, result interface{}) (int, error)
	Request(ctx context.Context, params Params, result interface{}) (int, error)
	PATCH(ctx context.Context, params Params, result interface{}) (int, error)
}

/*
Struct implements all the methods that are defined on the parent struct
It is useful to embed this struct in the struct that implements the parent struct
to act as a default method in case you want to implement only some of the methods and not all
*/
type UnImplemented struct{}

var _ Connector = (*UnImplemented)(nil)

func (unImpl UnImplemented) Request(
	context.Context,
	Params,
	interface{}) (int, error) {

	return 0, errUnimplemented
}
func (unImpl UnImplemented) GET(context.Context, Params, interface{}) (int, error) {
	return 0, errUnimplemented
}

func (unImpl UnImplemented) POST(context.Context, Params, interface{}) (int, error) {
	return 0, errUnimplemented
}

func (unImpl UnImplemented) PATCH(context.Context, Params, interface{}) (int, error) {
	return 0, errUnimplemented
}

// *************************** //
