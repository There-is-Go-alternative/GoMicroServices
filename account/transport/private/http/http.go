package http

import "net/http"

type AuthHTTP struct {
	http.Client
}

func (h AuthHTTP) ValidateToken(token string) bool {
	return true
}
