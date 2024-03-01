package apierror

import "fmt"

type ApiError struct {
	Code    int
	Message string
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
	//return string(e.Code) + ": " + e.Message
}

func New(code int, message string) *ApiError {
	return &ApiError{Code: code, Message: message}
}

func NewWithError(code int, message string, err error) *ApiError {
	return &ApiError{Code: code, Message: message + " / " + err.Error()}
}

// func New(code int, message string) *ApiError {
// 	return &ApiError{Code: code, Message: message}
// }
