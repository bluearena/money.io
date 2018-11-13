package errors

import (
	"net/http"
)

// bundle's error codes
const (
	ErrorCodeUnknown = 500
)

// ResponseError represents bundle errors structure
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewResponseError creates response error based on arguments
func NewResponseError(code int, message string) *ResponseError {
	return &ResponseError{
		Code:    code,
		Message: message,
	}
}

// GetHTTPStatus maps error code to http code
func (b *ResponseError) GetHTTPStatus() int {
	return http.StatusOK
}
