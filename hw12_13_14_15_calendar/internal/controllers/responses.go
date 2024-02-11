package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	ErrNo = iota
	ErrHas
)

var ErrEncodeJson = errors.New("error encode json")

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     int         `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Err     int    `json:"error"`
}

func sendErrorResponse(err error, w http.ResponseWriter) error {
	resp := ErrorResponse{
		Message: err.Error(),
		Err:     ErrHas,
	}
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return err
	}
	return nil
}
