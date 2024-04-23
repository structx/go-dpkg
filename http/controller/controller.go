// Package controller exposed http controllers
package controller

import "net/http"

// V0 exposed http controller on root handlers
type V0 interface {
	RegisterRoutesV0() http.Handler
}

// V1 exposed http controller on service handlers
type V1 interface {
	RegisterRoutesV1() http.Handler
}

// V1P exposed http controller for protected handlers
type V1P interface {
	RegisterRoutesV1P() http.Handler
}
