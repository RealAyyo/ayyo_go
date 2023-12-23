package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//nolint:exhaustive
func inValidate(inValues string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		has := false
		for _, val := range strings.Split(inValues, ",") {
			if v.String() == val {
				has = true
			}
		}
		if !has {
			return ErrInvalidIn
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		inVal := strings.Split(inValues, ",")
		has := false
		for _, val := range inVal {
			valInt, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			fmt.Println("-------")
			fmt.Println(v.Int(), "v.Int()")
			fmt.Println(int64(valInt), "int64(valInt)")
			fmt.Println("-------")
			if v.Int() == int64(valInt) {
				has = true
			}
		}
		if !has {
			return ErrInvalidIn
		}
	default:
		return ErrInvalidType
	}

	return nil
}
