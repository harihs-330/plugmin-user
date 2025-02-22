package errortools

import "fmt"

type errorCode string

type validationFails map[string][]errRecord

type errRecord struct {
	Code    errorCode `json:"code"`
	Message string    `json:"message"`
}

type Error struct {
	Code            errorCode       `json:"code"`
	Msg             string          `json:"msg"`
	Detail          string          `json:"details,omitempty"`
	HelpLink        string          `json:"helplink,omitempty"`
	ValidationFails validationFails `json:"validations,omitempty"`
}

func Init() *Error {
	return &Error{}
}

func (err *Error) Error() string {
	msg := "usr: " + err.Msg
	for k, rec := range err.ValidationFails {
		msg = fmt.Sprintf("%s { %v : ", msg, k)
		for _, v := range rec {
			msg += ": " + v.Message
		}
		msg += "}"
	}

	return msg
}

func (err *Error) Message() string {
	return err.Code.Name()
}

func (err *Error) AddValidationError(fieldName string, code errorCode, args ...any) {
	if err.Code == "" {
		err.Code = errorCode(ValidationFailure)
		err.Msg = err.Code.Name()
	}
	if err.ValidationFails == nil {
		err.ValidationFails = make(validationFails, 1)
	}
	msg := code.Name()
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	err.ValidationFails[fieldName] = append(err.ValidationFails[fieldName], errRecord{
		Code:    code,
		Message: msg,
	})
}

func (err *Error) Nil() bool {
	return err == nil || (err.Code == "" && len(err.ValidationFails) == 0)
}

func (erCode errorCode) Name() string {
	return errorCodeNames[string(erCode)]
}

func (erCode errorCode) HTTPCode() int {
	return errorHTTPCodes[string(erCode)]
}

// type Option
type Option interface {
	apply(*Error)
}

func (fails validationFails) apply(err *Error) {
	err.ValidationFails = fails
}

func WithValidationFails(fieldName string, code ...errorCode) Option {
	fails := make(validationFails, 1)
	for _, v := range code {
		fails[fieldName] = append(fails[fieldName], errRecord{
			Code:    v,
			Message: v.Name(),
		})
	}

	return Option(fails)
}

type DetailsOption string

func (detail DetailsOption) apply(err *Error) {
	err.Detail = string(detail)
	if err.Msg == "" {
		err.Msg = err.Detail
	}
}

func WithDetail(detail string) Option {
	return DetailsOption(detail)
}

func New(code errorCode, opts ...Option) *Error {
	errNew := &Error{
		Code:            code,
		Msg:             code.Name(),
		ValidationFails: make(validationFails),
	}
	for _, v := range opts {
		v.apply(errNew)
	}

	return errNew
}
