package hw09structvalidator

import (
	"reflect"
	"regexp"
)

//nolint:exhaustive
func regexValidate(regexpValue string, v reflect.Value) error {
	re := regexp.MustCompile(regexpValue)

	switch v.Kind() {
	case reflect.String:
		str := re.MatchString(v.String())
		if !str {
			return ErrInvalidRegexp
		}
	default:
		return ErrInvalidType
	}
	return nil
}
