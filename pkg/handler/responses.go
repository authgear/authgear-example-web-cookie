package handler

import (
	"net/http"
)

func UnknownErrorResponse(w http.ResponseWriter) {
	http.Error(w, "Unknown Error", http.StatusInternalServerError)
}

func NotFoundResponse(w http.ResponseWriter, msg string) {
	text := msg
	if text == "" {
		text = "Not found"
	}
	http.Error(w, text, http.StatusNotFound)
}
