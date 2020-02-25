package types

import (
	"errors"
)

// errors
var (
	ErrBasicAuth  = errors.New("basic auth required")
	ErrInvalidURL = errors.New("invalid url")
)
