package response

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details any    `json:"details,omitempty"`
}

func RespondError(w http.ResponseWriter, status int, err *Error) {
	RespondJSON(w, status, err)
}

func RespondStatus(w http.ResponseWriter, status int) {
	RespondJSON(w, status, Error{
		Message: http.StatusText(status),
	})
}
