// Package middleware http router middlewares
package middleware

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/structx/go-dpkg/domain"
)

// Auth middleware implementation
type Auth struct {
	log *zap.SugaredLogger
	mb  domain.MessageBroker
}

// NewAuth constructor
func NewAuth(logger *zap.Logger, broker domain.MessageBroker) *Auth {
	return &Auth{
		log: logger.Sugar().Named("AuthMiddleware"),
		mb:  broker,
	}
}

// Authenticate http middleware to verify wallet
func (a *Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var p domain.Permission
		switch r.Method {
		case http.MethodGet:
			p = domain.Read
		case http.MethodPost:
			p = domain.Write
		case http.MethodPut:
			p = domain.Write
		case http.MethodDelete:
			p = domain.Delete
		default:
			a.log.Errorf("received supported http method")
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		bb, err := json.Marshal(&domain.AccessControlEntry{
			Subject:    "",
			Resource:   "",
			Permission: p,
		})
		if err != nil {
			a.log.Errorf("failed to marshal access control entry %v", err)
		}

		m := domain.NewMsg(domain.VerifyUserAccess.String(), bb)
		rr, err := a.mb.RequestResponse(ctx, m)
		if err != nil {
			a.log.Errorf("failed to receive response from request response %v", err)
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		var en domain.EntryResponse
		err = json.Unmarshal(rr.GetPayload(), &en)
		if err != nil {
			a.log.Errorf("failed to unmarshal entry response %v", err)
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		if !en.Granted {
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
