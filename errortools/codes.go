package errortools

import (
	"net/http"
	"user/consts"
)

const (
	InternalServerErrCode = "USR500"
)

const (
	ValidationFailure       = "USR4000"
	BindingError            = "USR4001"
	MissingQueryParam       = "USR4002"
	MinLengthError          = "USR4003"
	MissingField            = "USR4004"
	BadFormat               = "USR4005"
	DataFetchFail           = "USR4006"
	JSONSyntaxError         = "USR4007"
	NoRecord                = "USR4008"
	InvalidFieldType        = "USR4009"
	ValidatorErrors         = "USR4010"
	DuplicateRecord         = "USR4011"
	InvalidToken            = "USR4012"
	InvalidPageErrorCode    = "USR4013"
	InvalidLimitErrorCode   = "USR4014"
	InvalidSearchKey        = "USR4015"
	InvalidSortOrder        = "USR4016"
	InvalidSortBy           = "USR4017"
	ProjectIDRequired       = "USR4018"
	InvalidUUID             = "USR4019"
	TokenExpired            = "USR4020"
	NoPermission            = "USR4021"
	FailedGettingPermission = "USR4022"
	MissingPathParam        = "USR4023"
	InvalidEmailID          = "USR4024"
	WrongPassword           = "USR4025"
)

const (
	UnauthorizedAccess = "USR600"
)

var errorCodeNames = map[string]string{
	InternalServerErrCode: "internal server error",
	MissingQueryParam:     "mandatory to pass the query param",
	ValidationFailure:     "Validation error",
	MinLengthError:        "minimum length should be %v characters",
	BindingError: `The request parameter was not bound correctly.
	please check the input params and try again`,
	MissingField:            "%s is mandatory.",
	BadFormat:               "Please check the format of the field",
	DataFetchFail:           "Couldn't get the  requested details. Please check the inputs",
	JSONSyntaxError:         "JSON syntax error at byte offset",
	NoRecord:                "No records found",
	UnauthorizedAccess:      "Unauthorized access",
	ValidatorErrors:         "%s",
	InvalidFieldType:        "Field is expecting a type of %s but got %s",
	DuplicateRecord:         "Record already exists",
	TokenExpired:            "Token has expired",
	InvalidToken:            "Token is invalid",
	InvalidPageErrorCode:    "Page must be a positive number greater than zero.",
	InvalidLimitErrorCode:   "Pagination limit should be between 1 and %v",
	InvalidSearchKey:        "Invalid search key %s",
	InvalidSortOrder:        "Sort order must be 'asc' or 'desc'.",
	InvalidSortBy:           "Sort value must be one of the following: %s",
	ProjectIDRequired:       "Project id is required",
	InvalidUUID:             "Invalid uuid value for %s",
	MissingPathParam:        "Missing Required Path Parameter",
	NoPermission:            "The user doesnot currently have any permissions on the project",
	FailedGettingPermission: "Error occurred,failed checking user-project relation",
	WrongPassword:           "The password you entered is incorrect",
	InvalidEmailID:          "The email address you entered is not registered",
}

var errorHTTPCodes = map[string]int{
	InternalServerErrCode: http.StatusInternalServerError,
	ValidationFailure:     http.StatusBadRequest,
	UnauthorizedAccess:    http.StatusUnauthorized,
	NoRecord:              http.StatusNotFound,
	TokenExpired:          consts.StatusInvalidToken,
}

var CustomMessages = map[string]string{
	"required":     "%s is a required field",
	"trimmedempty": "%s cannot be empty after trimming",
	"eqfield":      "Password and Confirm Password do not match",
	//nolint
	"passwordformat": "Password must be 8-20 characters long and include at least one uppercase letter, one lowercase letter, one digit, and one special character",
	"oneof":          "The '%s' field must be one of the following values: %s.",
	"min":            "The %s field must contain a minimum of %s characters",
	"max":            "The %s field must not exceed %s characters",
	//nolint
	"customname":  "The %s field should contains only alphabetic characters and spaces, with no leading, trailing, or consecutive spaces",
	"customemail": "Please provide a valid email in the format example@domain.com",
}
