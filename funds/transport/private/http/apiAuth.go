package http

import (
	"fmt"
	"net/http"
	"strings"
)

type AuthHTTP struct {
	http.Client
	AuthorizedKeys []string
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (h AuthHTTP) ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("no authorization header")
	}

	hasPrefix := strings.HasPrefix(token, "Bearer ")

	if !hasPrefix {
		return fmt.Errorf("not a bearer token")
	}

	sliced_token := token[7:]

	if !stringInSlice(sliced_token, h.AuthorizedKeys) {
		return fmt.Errorf("key not authorized")
	}
	return nil
}
