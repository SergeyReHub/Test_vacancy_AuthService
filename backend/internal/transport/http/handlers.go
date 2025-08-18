package http

import (
	"auth/backend/internal/service"
	"auth/backend/internal/transport/models"
	"encoding/json"
	"net/http"
	"strings"
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
// @Param Authorization header string true "Authorization header" default(Bearer <guid>)
// @Success      200  {object}  models.TokenResponse
// @Failure      500  {object}  models.ErrorResponse
// @Failure      401  {string}  string  "You haven't authorized and gotten guid"
// @Router       /auth/take_both_tokens [post]
func (h *Handlers) TakeBothTokens(w http.ResponseWriter, r *http.Request) {
	// 1. Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	userAgent := r.Header.Get("User-Agent")

	// 2. Check if it exists
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Missing Authorization header"}`))
		return
	}

	// 3. Get the token (simple version)
	guid := strings.Replace(authHeader, "Bearer ", "", 1)

	tokens, err := h.Service.TakeBothTokens(guid, userAgent, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	// 5. Return the tokens
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

// RefreshTokens godoc
// @Summary      Refresh both tokens
// @Description  Returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token header string true "Refresh token"
// @Param        access_token header string true "Access token"
// @Success      200  {object}  models.TokenResponse
// @Failure      500  {object}  models.ErrorResponse
// @Failure      403  {string}  string  "Refresh token family has been blocked"
// @Failure      401  {string}  string  "You have been deauthorized"
// @Failure      400  {string}  string  "Nil refresh or access tokens headers"
// @Failure      404  {string}  string  "Refresh token is not found"
// @Router       /auth/refresh_tokens [post]
func (h *Handlers) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("refresh_token")
	accessToken := r.Header.Get("access_token")
	userAgent := r.Header.Get("User-Agent")

	if refreshToken == "" || accessToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Nil refresh or access tokens headers"))
		return
	}

	tokens, err := h.Service.RefreshTokens(refreshToken, accessToken, userAgent, r.Context())
	if err != nil {
		switch err.Error() {
		case "Changed User-Agent":
			h.Deauthorization(w, r)
			return
		case "User deauthorized":
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("You have been deauthorized"))
			return
		case "Token is not valid. Now all token family is invalid":
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Refresh token family has been blocked"))
			return
		case "Refresh token not found":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Refresh token is not found"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("refresh_token", tokens.RefreshToken)
	w.Header().Set("access_token", tokens.AccessToken)
	w.WriteHeader(http.StatusOK)
}

// TakeGUID godoc
// @Summary      Get GUID
// @Description  Returns user GUID after successful authentication
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.User true "User credentials"
// @Success      200  {object}  models.GuidResponse
// @Failure      500  {object}  models.ErrorResponse
// @Failure      400  {string}  string  "Invalid request format"
// @Failure      401  {string}  string  "Invalid credentials"
// @Router       /auth/take_guid [post]
func (h *Handlers) TakeGUID(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request format"))
		return
	}

	if user.Username == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username and password are required"))
		return
	}

	guid, err := h.Service.TakeGUID(&user, r.Context())
	if err != nil {
		if err.Error() == "user deauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid credentials"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("guid", guid)
	w.WriteHeader(http.StatusOK)
}

// Deauthorization godoc
// @Summary      Deauthorize user
// @Description  Makes refresh token invalid
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        access_token header string true "Access token"
// @Success      200  {string}  string  "Deauthorized"
// @Failure      500  {object}  models.ErrorResponse
// @Failure      400  {string}  string  "Refresh token is not exists"
// @Router       /auth/deauthorization [post]
func (h *Handlers) Deauthorization(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("access_token")
	if accessToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Access token is required"))
		return
	}

	err := h.Service.Deauthorization(accessToken, r.Context())
	if err != nil {
		if err.Error() == "Refresh token is not exists" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Refresh token is not exists"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deauthorized"))
}
