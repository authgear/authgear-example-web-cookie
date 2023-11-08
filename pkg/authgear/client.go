package authgear

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type Client struct {
	Endpoint *url.URL
}

func NewClient(endpoint string) *Client {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic("invalid endpoint")
	}

	return &Client{
		Endpoint: u,
	}
}

func (h *Client) GetUserInfo(cookie string) (map[string]interface{}, error) {

	url := h.Endpoint.JoinPath("/oauth2/userinfo")

	httpReq, err := http.NewRequest("GET", url.String(), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Cookie", cookie)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
