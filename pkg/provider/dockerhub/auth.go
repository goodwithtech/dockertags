package dockerhub

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
)

func getJSON(ctx context.Context, url string, auth types.AuthConfig, timeout time.Duration, response interface{}) (http.Header, error) {
	cli, err := new(auth, timeout)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return resp.Header, nil
}

func new(auth types.AuthConfig, timeout time.Duration) (*http.Client, error) {
	transport := http.DefaultTransport
	tokenTransport := &tokenTransport{
		Transport: transport,
		Username:  auth.Username,
		Password:  auth.Password,
	}

	registry := &http.Client{
		Timeout:   timeout,
		Transport: tokenTransport,
	}
	return registry, nil
}
