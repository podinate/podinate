package apierror

import (
	"fmt"
	"net/http"

	api "github.com/johncave/podinate/controller/go"
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

func (e ApiError) EncodeJSONResponse(w http.ResponseWriter) {

	//out, _ := json.MarshalIndent(e, "", "  ")
	//resp := api.ImplResponse{Code: e.Code, Body: string(out)}
	api.EncodeJSONResponse(e, &e.Code, w)
}

// func New(code int, message string) *ApiError {
// 	return &ApiError{Code: code, Message: message}
// }
