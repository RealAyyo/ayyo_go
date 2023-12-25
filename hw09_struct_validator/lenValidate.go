package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

//nolint:exhaustive
func lenValidate(requireLen string, v reflect.Value) error {
	maxLen, err := strconv.Atoi(requireLen)
	if err != nil {
		return fmt.Errorf("invalid len argument")
	}

	switch v.Kind() {
	case reflect.String:
		if len(v.String()) != maxLen {
			return ErrInvalidLen
		}
	case reflect.Slice:
		for _, val := range v.Interface().([]string) {
			if len(val) != maxLen {
				return ErrInvalidLen
			}
		}
	default:
		return ErrInvalidType
	}

	return nil
}
