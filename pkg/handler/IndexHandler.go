package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/oursky/authgear-exmaple-web-cookie/pkg/authgear"
)

type IndexHandler struct {
	AuthgearClient   *authgear.Client
	AuthgearEndpoint string
}

var _ http.Handler = &IndexHandler{}

type Data struct {
	AuthgearEndpoint string
	IsAuthenticated  bool
	UserJSON         string
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles("templates/index.html", "templates/index.css.html")
	if err != nil {
		UnknownErrorResponse(w)
		return
	}

	userID := r.Header.Get("X-Authgear-User-Id")

	var user map[string]interface{}

	user, err = h.AuthgearClient.GetUserInfo(r.Header.Get("Cookie"))
	if err != nil {
		user = map[string]interface{}{
			"error": err.Error(),
		}
	}

	userJSON, _ := json.MarshalIndent(user, "", "  ")

	data := &Data{
		IsAuthenticated:  userID != "",
		AuthgearEndpoint: h.AuthgearEndpoint,
		UserJSON:         string(userJSON),
	}

	err = t.Execute(w, data)
	if err != nil {
		UnknownErrorResponse(w)
		return
	}
}
