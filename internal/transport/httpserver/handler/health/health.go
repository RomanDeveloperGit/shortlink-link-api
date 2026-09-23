package health

import "net/http"

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
