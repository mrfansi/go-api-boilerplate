package errors

import "errors"

var (
	// Common errors
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrConflict          = errors.New("resource already exists")
	ErrValidation        = errors.New("validation error")
	ErrBadRequest        = errors.New("bad request")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidCredential = errors.New("invalid credentials")

	// User specific errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidRole       = errors.New("invalid role")
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}
