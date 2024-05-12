// Package controller exposed http controllers
package controller

import (
	"github.com/go-chi/chi/v5"
)

// V0 exposed http controller on root handlers
type V0 interface {
	RegisterRoutesV0(r chi.Router)
}

// V1 exposed http controller on service handlers
type V1 interface {
	RegisterRoutesV1(r chi.Router)
}

// V1P exposed http controller for protected handlers
type V1P interface {
	RegisterRoutesV1P(r chi.Router)
}
