package utils

import "net/http"

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AppError struct {
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Helpers
func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Status:  http.StatusNotFound,
		Message: message,
		Err:     err,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Status:  http.StatusUnauthorized,
		Message: message,
		Err:     err,
	}
}

func NewInternal(err error) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Message: "Erro interno do servidor",
		Err:     err,
	}
}
