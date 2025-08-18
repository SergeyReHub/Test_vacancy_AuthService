package http

import (
	"auth/backend/internal/service"
	_ "auth/docs"

	"context"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, service service.AuthServiceUC, ctx context.Context) *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	handlers := NewHandlers(service)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/take_both_tokens", handlers.TakeBothTokens).Methods("POST")
	auth.HandleFunc("/refresh_tokens", handlers.RefreshTokens).Methods("POST")
	auth.HandleFunc("/take_guid", handlers.TakeGUID).Methods("POST")
	auth.HandleFunc("/deauthorization", handlers.Deauthorization).Methods("POST")

	return r
}
