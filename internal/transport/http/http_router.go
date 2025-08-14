package http

import (
	"auth/internal/service"
	"context"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, service service.AuthServiceUC, ctx context.Context) *mux.Router {
	r := mux.NewRouter()

	handlers := NewHandlers(service)

	auth := r.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/take_both_tokens", handlers.TakeBothTokens).Methods().Headers("Authorization", "Bearer.*")
	auth.HandleFunc("/refresh_tokens", handlers.RefreshTokens)
	auth.HandleFunc("/take_guid", handlers.TakeGUID)
	auth.HandleFunc("/deauthorization", handlers.Deauthorization)

	return r
}
