package validators

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var (
	ErrBadRequest = errors.New("bad Request")
	Validator     = validator.New()
)

func Validate(method string, r *http.Request, data interface{}) error {
	if method != r.Method {
		return ErrBadRequest
	}

	switch method {
	case "POST":
	case "PATCH":
		return ValidatePostQuery(r, data)
	case "GET":
		return ValidateGetQuery(r, data)
	}

	return ErrBadRequest
}

func ValidatePostQuery(r *http.Request, data interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrBadRequest
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return ErrBadRequest
	}

	err = Validator.Struct(data)
	if err != nil {
		return ErrBadRequest
	}

	return nil
}

func ValidateGetQuery(r *http.Request, data interface{}) error {
	queryParams := r.URL.Query()

	queryData := make(map[string]interface{})
	for key, values := range queryParams {
		queryData[key] = values[0]
	}

	jsonData, err := json.Marshal(queryData)
	if err != nil {
		return ErrBadRequest
	}

	err = json.Unmarshal(jsonData, data)
	if err != nil {
		return ErrBadRequest
	}

	err = Validator.Struct(data)
	if err != nil {
		return ErrBadRequest
	}

	return nil
}
