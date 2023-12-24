package hw09structvalidator

import (
	"reflect"
	"strconv"
	"strings"
)

//nolint:exhaustive
func inValidate(inValues string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		for _, val := range strings.Split(inValues, ",") {
			if v.String() == val {
				return nil
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		inVal := strings.Split(inValues, ",")

		for _, val := range inVal {
			valInt, err := strconv.Atoi(val)
			if err != nil {
				return err
			}

			if v.Int() == int64(valInt) {
				return nil
			}
		}

	default:
		return ErrInvalidType
	}

	return ErrInvalidIn
}
