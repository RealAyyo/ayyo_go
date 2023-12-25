package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrNotStruct     = errors.New("input value is not a struct")
	ErrInvalidLen    = errors.New("invalid length")
	ErrInvalidMax    = errors.New("invalid max value")
	ErrInvalidIn     = errors.New("invalid in value")
	ErrInvalidMin    = errors.New("invalid min value")
	ErrInvalidType   = errors.New("invalid type")
	ErrInvalidRegexp = errors.New("invalid regexp value")
)

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("validation failed: ")
	for i, err := range v {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Err))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var errs ValidationErrors
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)

		tag := rt.Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		if validationErrs := validateField(field, tag); validationErrs != nil {
			for _, err := range validationErrs {
				errs = append(errs, ValidationError{
					Field: rt.Field(i).Name,
					Err:   err,
				})
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateField(v reflect.Value, tag string) []error {
	validateRules := strings.Split(tag, "|")
	var errs []error

	for _, ruleField := range validateRules {
		if len(ruleField) < 2 {
			continue
		}

		args := strings.Split(ruleField, ":")[1:][0]
		rule := strings.Split(ruleField, ":")[0]

		switch rule {
		case "len":
			err := lenValidate(args, v)
			if err != nil {
				errs = append(errs, err)
			}
		case "regexp":
			err := regexValidate(args, v)
			if err != nil {
				errs = append(errs, err)
			}
		case "in":
			err := inValidate(args, v)
			if err != nil {
				errs = append(errs, err)
			}
		case "min":
			err := minValidate(args, v)
			if err != nil {
				errs = append(errs, err)
			}
		case "max":
			err := maxValidate(args, v)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
