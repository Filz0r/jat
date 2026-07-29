package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) > 0 {
		split := strings.Split(auth, " ")
		if len(split) == 2 && split[0] == "Bearer" {
			return split[1], nil
		}
	}

	if c, err := r.Cookie("access_token"); err == nil && c.Value != "" {
		return c.Value, nil
	}

	return "", fmt.Errorf("no bearer token found")
}

func GetAPIKey(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) > 0 {
		split := strings.Split(auth, " ")
		if len(split) == 2 && split[0] == "Bearer" {
			return split[1], nil
		}
	}
	if c, err := r.Cookie("refresh_token"); err == nil && c.Value != "" {
		return c.Value, nil
	}

	return "", fmt.Errorf("no refresh/api key found")
}
