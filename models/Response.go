package models

type Response struct {
	ResponseCode uint16 `json:"responsecode"`
	Message      string `json:"message"`
}

type ErrorResponseRegisterLogin struct {
	ResponseCode int    `json:"responsecode"`
	Message      string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}