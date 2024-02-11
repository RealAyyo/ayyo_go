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
	Validate      = validator.New()
)

type Validator struct {
}

func NewQueryValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(method string, r *http.Request, data interface{}) error {
	if method != r.Method {
		return ErrBadRequest
	}

	switch method {
	case "POST":
		return v.ValidatePostQuery(r, data)
	case "PATCH":
		return v.ValidatePostQuery(r, data)
	case "GET":
		return v.ValidateGetQuery(r, data)
	}

	return ErrBadRequest
}

func (v *Validator) ValidatePostQuery(r *http.Request, data interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrBadRequest
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return ErrBadRequest
	}

	err = Validate.Struct(data)
	if err != nil {
		return ErrBadRequest
	}

	return nil
}

func (v *Validator) ValidateGetQuery(r *http.Request, data interface{}) error {
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

	err = Validate.Struct(data)
	if err != nil {
		return ErrBadRequest
	}

	return nil
}
