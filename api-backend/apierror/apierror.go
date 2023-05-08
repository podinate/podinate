package apierror

type ApiError struct {
	Code    int
	Message string
}

func (e ApiError) Error() string {
	return string(e.Code) + ": " + e.Message
}

func New(code int, message string) ApiError {
	return ApiError{Code: code, Message: message}
}

// func New(code int, message string) *ApiError {
// 	return &ApiError{Code: code, Message: message}
// }
