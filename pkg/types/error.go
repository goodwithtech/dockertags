package types

import (
	"errors"
	"time"
)

var (
	ErrBasicAuth  = errors.New("basic auth required")
	ErrInvalidURL = errors.New("invalid url")
)

type AuthOption struct {
	AuthURL      string
	UserName     string
	Password     string
	GcpCredPath  string
	AwsAccessKey string
	AwsSecretKey string
	AwsRegion    string
	Timeout      time.Duration
}
