package structutil

import "reflect"

// GetStructFieldCount returns the number of fields in the provided struct.
// It panics if the provided value is not a struct.
func GetStructFieldCount(v interface{}) int {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		panic("provided value is not a struct")
	}

	return val.Type().NumField()
}
