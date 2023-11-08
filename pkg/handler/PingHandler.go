package handler

import (
	"encoding/json"
	"net/http"
)

type PingHandler struct {
}

var _ http.Handler = &PingHandler{}

type PingResponse struct {
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(&PingResponse{})
	if err != nil {
		UnknownErrorResponse(w)
		return
	}
}
