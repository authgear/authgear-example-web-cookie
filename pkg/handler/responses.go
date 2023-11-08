package handler

import "net/http"

func UnknownErrorResponse(w http.ResponseWriter) {
	http.Error(w, "Unknown Error", http.StatusInternalServerError)
}
