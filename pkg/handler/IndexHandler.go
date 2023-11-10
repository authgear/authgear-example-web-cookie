package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/authgear/authgear-exmaple-web-cookie/pkg/authgear"
	"github.com/authgear/authgear-server/pkg/util/log"
)

type IndexHandler struct {
	Logger           *log.Logger
	AuthgearClient   *authgear.Client
	AuthgearEndpoint *url.URL
	DefaultClientID  string
}

var _ http.Handler = &IndexHandler{}

type Data struct {
	AuthgearEndpoint string
	IsAuthenticated  bool
	UserJSON         string
	DefaultClientID  string
}

const (
	CookieNameDefaultClientID string = "default_client_id"
)

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	oidcConfig, err := h.AuthgearClient.GetOIDCConfig()
	if err != nil {
		h.Logger.Errorln("failed to get oidc config")
		UnknownErrorResponse(w)
		return
	}
	userInfoEndpoint := oidcConfig["userinfo_endpoint"].(string)
	authorizationEndpoint := oidcConfig["authorization_endpoint"].(string)
	endSessionEndpoint := oidcConfig["end_session_endpoint"].(string)

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		t, err := template.ParseFiles("templates/index.html", "templates/index.css.html")
		if err != nil {
			h.Logger.Errorln("failed to parse templates")
			UnknownErrorResponse(w)
			return
		}

		userID := r.Header.Get("X-Authgear-User-Id")

		var user map[string]interface{}

		user, err = h.AuthgearClient.GetUserInfo(userInfoEndpoint, r.Header.Get("Cookie"))
		if err != nil {
			user = map[string]interface{}{
				"error": err.Error(),
			}
		}

		userJSON, _ := json.MarshalIndent(user, "", "  ")

		defaultCientID := h.DefaultClientID

		defaultCientIDCookie, err := r.Cookie(CookieNameDefaultClientID)
		if err == nil && defaultCientIDCookie.Value != "" {
			defaultCientID = defaultCientIDCookie.Value
		}

		data := &Data{
			IsAuthenticated:  userID != "",
			AuthgearEndpoint: h.AuthgearEndpoint.String(),
			UserJSON:         string(userJSON),
			DefaultClientID:  defaultCientID,
		}

		err = t.Execute(w, data)
		if err != nil {
			h.Logger.Errorln("failed to execute templates")
			UnknownErrorResponse(w)
			return
		}
		return

	case "POST":
		err := r.ParseForm()
		if err != nil {
			h.Logger.Errorln("failed to parse form")
			UnknownErrorResponse(w)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "login":
			clientID := r.FormValue("client_id")
			q := url.Values{}
			redirectURI := r.Header.Get("Origin") + "/"
			q.Set("redirect_uri", redirectURI)
			q.Set("scope", strings.Join([]string{
				"openid", "offline_access", "https://authgear.com/scopes/full-access",
			}, " "))
			q.Set("response_type", "none")
			q.Set("client_id", clientID)
			q.Set("x_sso_enabled", "true")
			q.Set("prompt", "login")

			url, err := url.Parse(authorizationEndpoint)
			if err != nil {
				h.Logger.Errorln("failed to parse authorizationEndpoint")
				UnknownErrorResponse(w)
				return
			}

			url.RawQuery = q.Encode()

			http.SetCookie(w, &http.Cookie{Name: CookieNameDefaultClientID, Value: clientID})

			http.Redirect(w, r, url.String(), http.StatusFound)
			return
		case "logout":
			logoutEndpoint, err := url.Parse(endSessionEndpoint)
			if err != nil {
				h.Logger.Errorln("failed to parse endSessionEndpoint")
				UnknownErrorResponse(w)
				return
			}
			redirectURI := r.Header.Get("Origin") + "/"
			q := url.Values{}
			q.Set("post_logout_redirect_uri", redirectURI)
			url := logoutEndpoint
			url.RawQuery = q.Encode()
			http.Redirect(w, r, url.String(), http.StatusFound)
			return
		default:
			NotFoundResponse(w, "")
			return
		}

	default:
		NotFoundResponse(w, "")
		return
	}
}
