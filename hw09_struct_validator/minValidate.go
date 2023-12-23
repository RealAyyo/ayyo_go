package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

func minValidate(requireMax string, v reflect.Value) error {
	requireMaxVal, err := strconv.Atoi(requireMax)
	if err != nil {
		return fmt.Errorf("invalid Max argument")
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < int64(requireMaxVal) {
			return ErrInvalidMin
		}
	case reflect.Slice:
		for _, val := range v.Interface().([]int) {
			if val < requireMaxVal {
				return ErrInvalidMin
			}
		}
	default:
		return ErrInvalidType
	}
	return nil
}
