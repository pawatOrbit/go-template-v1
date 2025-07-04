package validation

import (
	"errors"
	"reflect"
	"strings"
)

// ValidateStruct checks for fields with the `validate:"true"` tag and ensures they are not empty.
func ValidateStruct(s any) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("expected a struct or a pointer to a struct")
	}

	errorFields := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := val.Type().Field(i).Tag.Get("validate")
		if tag == "true" && field.IsZero() {
			errorFields = append(errorFields, val.Type().Field(i).Name)
		}
	}

	if len(errorFields) > 0 {
		return errors.New("fields are empty or zero value: " + strings.Join(errorFields, ", "))
	}

	return nil
}
