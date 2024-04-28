package domain

import "net/http"

// Authenticator http auth middleware interface
//
//go:generate mockery --name Authenticator
type Authenticator interface {
	// Authenticate http request
	Authenticate(http.Handler) http.Handler
}
