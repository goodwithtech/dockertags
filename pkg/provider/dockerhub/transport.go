package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const authURL = "https://hub.docker.com/v2/users/login/"

type authToken struct {
	Token string `json:"token"`
}

type tokenTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
}

// RoundTrip defines the round tripper for token transport.
func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	if t.Username == "" || t.Password == "" {
		return resp, nil
	}
	if !isTokenDemand(resp) {
		return resp, nil
	}
	resp.Body.Close()
	return t.authAndRetry(req)
}

func isTokenDemand(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	if resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}

func (t *tokenTransport) authAndRetry(req *http.Request) (*http.Response, error) {
	token, authResp, err := t.auth(req.Context())
	if err != nil {
		return authResp, err
	}

	response, err := t.retry(req, token)
	if response != nil {
		response.Header.Set("request-token", token)
	}
	return response, err
}

func (t *tokenTransport) retry(req *http.Request, token string) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return t.Transport.RoundTrip(req)
}

func (t *tokenTransport) auth(ctx context.Context) (string, *http.Response, error) {
	jsonStr := []byte(fmt.Sprintf(`{"username": "%s","password": "%s"}`, t.Username, t.Password))
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	c := http.Client{
		Transport: t.Transport,
	}
	resp, err := c.Do(req.WithContext(ctx))
	//resp, err := http.Post(authURL, "application/json", )

	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp, err
	}
	var authToken authToken
	if err := json.NewDecoder(resp.Body).Decode(&authToken); err != nil {
		return "", nil, err
	}
	token := authToken.Token
	return token, nil, nil
}
