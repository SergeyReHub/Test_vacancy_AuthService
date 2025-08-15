package http

import (
	"auth/internal/service"
	"auth/internal/transport/models"
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
		w.WriteHeader(http.StatusBadRequest)
	}
	tokens, err := h.Service.RefreshTokens(refresh_token, access_token, userAgent, r.Context())
	if err != nil {
		if err.Error() == "Changed User-Agent" {
			h.Deauthorization(w, r) // serve deauthorize
			return
		}
		if err.Error() == "User deauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err.Error() == "Token is not valid. Now all token family is invalid" {
			w.WriteHeader(http.StatusForbidden)
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
	user_name := r.Header.Get("user_name")
	user_password := r.Header.Get("user_password")

	if user_name == "" || user_password == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	guid, err := h.Service.TakeGUID(&models.User{
		Username: user_name,
		Password: user_password,
	}, r.Context())
	if err != nil {
		if err.Error() == "user deauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("guid", guid)
	w.WriteHeader(http.StatusOK)
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
	access_token := r.Header.Get("access_token")

	err := h.Service.Deauthorization(access_token, r.Context())
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
