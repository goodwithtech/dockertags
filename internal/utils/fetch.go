package utils

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/goodwithtech/dockertags/internal/log"

	"github.com/parnurzeal/gorequest"
)

var versionPattern = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)

func fetchURL(url string, cookie *http.Cookie) ([]byte, error) {
	resp, body, err := gorequest.New().AddCookie(cookie).Get(url).Type("text").EndBytes()
	if err != nil {
		return nil, fmt.Errorf("fail to fetch : %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error code : %d, url : %s", resp.StatusCode, url)
	}
	return body, nil
}

// FetchLatestVersion returns latest dockertags version
func FetchLatestVersion() (version string, err error) {
	log.Logger.Debug("Fetch latest version from github")
	body, err := fetchURL(
		"https://github.com/goodwithtech/dockertags/releases/latest",
		&http.Cookie{Name: "user_session", Value: "tags"},
	)
	if err != nil {
		return "", err
	}
	versionMatched := versionPattern.FindString(string(body))
	return versionMatched, nil
}
