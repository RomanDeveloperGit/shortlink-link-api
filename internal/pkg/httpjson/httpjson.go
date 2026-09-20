package httpjson

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, code string, msg string) {
	WriteJSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: msg,
		},
	})
}
