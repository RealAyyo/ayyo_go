package controllers

const (
	ErrNo = iota
	ErrHas
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     int         `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Err     int    `json:"error"`
}
