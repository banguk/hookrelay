package validate

import (
	"errors"
	"strings"
)

var (
	ErrEmptyID   = errors.New("id is empty")
	ErrEmptyURL  = errors.New("url is empty")
	ErrBadScheme = errors.New("url must start with http:// or https://")
)

func Webhook(id, url string) error {
	if id == "" {
		return ErrEmptyID
	}
	if url == "" {
		return ErrEmptyURL
	}
	if !strings.HasPrefix(url, "http://") {
		return ErrBadScheme
	}
	return nil
}
