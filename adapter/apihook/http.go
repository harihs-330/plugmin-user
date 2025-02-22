package apihook

import (
	"context"
	"net/http"
	"user/utils/apikit"
)

type HTTPAPI struct {
}

var _ Connector = (*HTTPAPI)(nil)

func (h HTTPAPI) Request(ctx context.Context, params Params, result interface{}) (int, error) {
	return apikit.APIRequest(
		ctx,
		params.URL, params.Method,
		params.Header, params.Body,
		result)
}

func (h HTTPAPI) GET(ctx context.Context, params Params, result interface{}) (int, error) {
	return apikit.APIRequest(
		ctx,
		params.URL, http.MethodGet,
		params.Header, params.Body,
		result)
}

func (h HTTPAPI) POST(ctx context.Context, params Params, result interface{}) (int, error) {
	return apikit.APIRequest(
		ctx,
		params.URL, http.MethodPost,
		params.Header, params.Body,
		result)
}

func (h HTTPAPI) PATCH(ctx context.Context, params Params, result interface{}) (int, error) {
	return apikit.APIRequest(
		ctx,
		params.URL, http.MethodPatch,
		params.Header, params.Body,
		result,
	)
}

func (h HTTPAPI) PUT(ctx context.Context, params Params, result interface{}) (int, error) {
	return apikit.APIRequest(
		ctx,
		params.URL, http.MethodPut,
		params.Header, params.Body,
		result,
	)
}
