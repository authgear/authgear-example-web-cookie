package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type IndexHandler struct {
	AuthgearEndpoint string
}

var _ http.Handler = &IndexHandler{}

type Data struct {
	AuthgearEndpoint string
	IsAuthenticated  bool
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles("templates/index.html", "templates/index.css.html")
	if err != nil {
		UnknownErrorResponse(w)
		return
	}
	fmt.Println("Headers", r.Header)

	userID := r.Header.Get("X-Authgear-User-Id")

	data := &Data{
		IsAuthenticated: userID != "",
	}

	err = t.Execute(w, data)
	if err != nil {
		UnknownErrorResponse(w)
		return
	}
}
