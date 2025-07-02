package utils

import "net/http"

type AppError struct {
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Message: message,
		Status:  http.StatusConflict,
	}
}

func NewInternalError() *AppError {
	return &AppError{
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}
}
