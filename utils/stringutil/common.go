package stringutil

import (
	"fmt"
	"math/big"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"crypto/rand"

	"github.com/google/uuid"
)

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func OTP(length int) int {
	otp := make([]string, length)
	for index := range otp {
		val, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			val = big.NewInt(1)
		}
		if index == 0 && val.Int64() == 0 {
			val, err = rand.Int(rand.Reader, big.NewInt(9))
			if err != nil {
				val = big.NewInt(1)
			}
			val = val.Add(val, big.NewInt(1))
		}
		otp[index] = strconv.Itoa(int(val.Int64()))
	}
	res, _ := strconv.Atoi(strings.Join(otp, ""))

	return res
}

// Function to create URL query parameters from a struct
func QueryParams(data interface{}) (string, error) {
	values := url.Values{}

	val, t, err := validateStruct(data)
	if err != nil {
		return "", err
	}

	for i := 0; i < val.NumField(); i++ {
		if tag, ok := getTag(t.Field(i), "url"); ok {
			if value, ok := getValue(val.Field(i)); ok {
				values.Set(tag, value)
			}
		}
	}

	return values.Encode(), nil
}

func validateStruct(data interface{}) (reflect.Value, reflect.Type, error) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("input is not a struct")
	}

	return v, t, nil
}

func getTag(field reflect.StructField, tagName string) (string, bool) {
	tag := field.Tag.Get(tagName)
	if tag == "" {
		return "", false
	}

	return tag, true
}

func getValue(field reflect.Value) (string, bool) {
	if !field.IsValid() {
		return "", false
	}

	if isZeroValue(field) {
		return "", false
	}

	return fmt.Sprintf("%v", field.Interface()), true
}

// isZeroValue checks if the given field is a zero value.
func isZeroValue(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.String:
		return field.String() == ""
	case reflect.Bool:
		return !field.Bool()
	default:
		return isZeroUUID(field)
	}
}

func isZeroUUID(field reflect.Value) bool {
	if field.Type() == reflect.TypeOf(uuid.UUID{}) {
		return field.Interface() == uuid.Nil
	}
	return false
}

// Function to check if an interface{} value is empty
func IsValueEmpty(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsZero()
}

type QueryParam string

func (qp QueryParam) ToString() string {
	return fmt.Sprintf("%v", qp)
}

func (qp QueryParam) ToInt() int {
	if value, err := strconv.Atoi(qp.ToString()); err == nil && value > 0 {
		return value
	}
	return -1
}

func ParseQueryParam(param string, defaultValue any) QueryParam {
	if IsEmpty(param) {
		return QueryParam(fmt.Sprintf("%v", defaultValue))
	}

	return QueryParam(param)
}

func RemoveDuplicates[T comparable](input []T) []T {
	uniqueMap := make(map[T]bool)
	uniqueList := []T{}

	for _, item := range input {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}
