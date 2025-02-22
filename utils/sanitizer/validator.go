package sanitizer

import (
	"net/mail"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Validator initializes and returns a new validator instance with custom validations
func Validator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(withJSONTag())

	// Register custom validation functions
	if err := validate.RegisterValidation("trimmedempty", isEmptyAfterTrim); err != nil {
		return validate
	}
	if err := validate.RegisterValidation("eqfield", validatePasswordConfirmation); err != nil {
		return validate
	}
	if err := validate.RegisterValidation("passwordformat", validatePasswordComplexity); err != nil {
		return validate
	}
	if err := validate.RegisterValidation("customemail", validateEmail); err != nil {
		return validate
	}
	if err := validate.RegisterValidation("customname", validateName); err != nil {
		return validate
	}

	return validate
}

// Custom validation function to check if a trimmed string is empty
func isEmptyAfterTrim(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

// Custom validation function to check if two fields are equal
func validatePasswordConfirmation(fl validator.FieldLevel) bool {
	// The field value
	fieldValue := fl.Field().String()

	// The value of the field to compare against
	fieldToCompare := fl.Parent().FieldByName(fl.Param()).String()

	return fieldValue == fieldToCompare
}

// Custom validation function to validate password complexity
func validatePasswordComplexity(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 8 && len(password) <= 20 &&
		containsUppercase(password) && containsLowercase(password) &&
		containsDigit(password) && containsSpecialChar(password)
}

// Check if a string contains at least one uppercase letter
func containsUppercase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}

	return false
}

// Check if a string contains at least one lowercase letter
func containsLowercase(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) {
			return true
		}
	}

	return false
}

// Check if a string contains at least one digit
func containsDigit(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}

	return false
}

// Check if a string contains at least one special character
func containsSpecialChar(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c) {
			return true
		}
	}

	return false
}

func validateEmail(fl validator.FieldLevel) bool {
	_, err := mail.ParseAddress(fl.Field().String())
	return err == nil
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	re := regexp.MustCompile(`^[a-zA-Z]+(?: [a-zA-Z]+)*$`)

	return re.MatchString(name)
}

// withJSONTag returns a function that extracts the JSON tag from struct fields
func withJSONTag() func(fld reflect.StructField) string {
	return func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// Skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	}
}
