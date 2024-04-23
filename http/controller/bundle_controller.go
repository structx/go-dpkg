package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Bundle controller
type Bundle struct {
	log *zap.SugaredLogger
}

// interface compliance
var _ V0 = (*Bundle)(nil)

// NewBundle constructor
func NewBundle(logger *zap.Logger) *Bundle {
	return &Bundle{
		log: logger.Sugar().Named("BundleController"),
	}
}

// RegisterRoutesV0 create handler from exposed routes
func (b *Bundle) RegisterRoutesV0() http.Handler {

	r := chi.NewRouter()

	r.Get("/health", b.Health)

	return r
}

// Health check handler
func (b *Bundle) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		b.log.Errorf("failed to write response %v", err)
		http.Error(w, "unable to write response", http.StatusInternalServerError)
	}
}
