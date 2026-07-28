package utils

import (
	"fmt"
	"net/url"
)

func ValidateDbUri(s string) error {
	if s == "" {
		return fmt.Errorf("database URL is required")
	}
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid database URL: %v", err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("invalid database URL: missing scheme")
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return fmt.Errorf("invalid database URL: invalid scheme")
	}
	return nil
}

func ValidateUrl(s string) error {
	if s == "" {
		return fmt.Errorf("a valid URL is required")
	}
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("invalid URL: missing scheme")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid URL: invalid scheme")
	}
	return nil
}
