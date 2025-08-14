package http

import (
	"auth/internal/service"
	"net/http"
)

type Handlers struct {
	Service service.AuthServiceUC
}

func NewHandlers(src service.AuthServiceUC) *Handlers {
	return &Handlers{
		Service: src,
	}
}

// TakeBothTokens godoc
// @Summary      Get both tokens
// @Description  Returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  TokenResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /auth/take_both_tokens [post]
func (h *Handlers) TakeBothTokens(w http.ResponseWriter, r *http.Request) {
	guid := r.Header.Get("guid")
	if guid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokens, err := h.Service.TakeBothTokens(guid, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("refresh_token", tokens.RefreshToken)
	w.Header().Add("access_token", tokens.AccessToken)
	w.WriteHeader(http.StatusOK)
}

// RefreshTokens godoc
// @Summary      Get both tokens
// @Description  Returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  TokenResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /auth/take_both_tokens [post]
func (h *Handlers) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	refresh_token := r.Header.Get("refresh_token")
	access_token := r.Header.Get("access_token")
	userAgent := r.Header.Get("User-Agent")

	if refresh_token == "" || access_token == "" {
		w.WriteHeader(http.StatusUnauthorized)
	}
	tokens, err := h.Service.RefreshTokens(refresh_token, access_token, userAgent, r.Context())
	if err != nil {
		if err.Error() == "Changed User-Agent" {
			http.Ser(w, h.Deauthorization) // serve deauthorize
		}
		if err.Error() == "User deauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err.Error() == "Token is not valid. Now all token family is invalid" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.Header().Add("refresh_token", tokens.RefreshToken)
	w.Header().Add("access_token", tokens.AccessToken)
}

// TakeGUID godoc
// @Summary      Get both tokens
// @Description  Returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  TokenResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /auth/take_both_tokens [post]
func (h *Handlers) TakeGUID(w http.ResponseWriter, r *http.Request) {

}

// Deauthorization godoc
// @Summary      Get both tokens
// @Description  Returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  TokenResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /auth/take_both_tokens [post]
func (h *Handlers) Deauthorization(w http.ResponseWriter, r *http.Request) {

}
