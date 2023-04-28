package responder

import (
	"encoding/json"
	"net/http"

	api "github.com/johncave/podinate/api-backend/go"
)

func Response(code int, body interface{}) api.ImplResponse {

	switch body.(type) {
	case error:
		out, _ := json.MarshalIndent(api.Model500Error{Code: float32(code), Message: body.(error).Error()}, "", "  ")
		return api.ImplResponse{Code: code, Body: string(out)}
	case []byte:
		return api.ImplResponse{Code: code, Body: string(body.([]byte))}
	case string:
		if code >= http.StatusBadRequest {
			return api.ImplResponse{Code: code, Body: api.Model500Error{Code: float32(code), Message: body.(string)}}
		}
		return api.ImplResponse{Code: code, Body: body.(string)}
	}
	return api.ImplResponse{Code: code, Body: body}
}
