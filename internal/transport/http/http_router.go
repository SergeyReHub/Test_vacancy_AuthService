package http

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, ctx context.Context) {
	r := mux.NewRouter()

	r.HandleFunc("/take_both_tokens", TakeBothTokens).Methods().Headers("Authorization", "Bearer.*")
	r.HandleFunc("/refresh_tokens", RefreshTokens)
	r.HandleFunc("/take_guid", TakeGUID)
	r.HandleFunc("/deauthorization", Deauthorization)

}
