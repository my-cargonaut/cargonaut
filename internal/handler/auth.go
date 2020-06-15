package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/my-cargonaut/cargonaut"
	"github.com/my-cargonaut/cargonaut/internal/jwt"
	"github.com/my-cargonaut/cargonaut/pkg/password"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token       string    `json:"token"`
	TokenExpiry time.Time `json:"token_expiry"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Get user identified by his unique email. If the user is not found, the
	// method does not return immediately. Instead some time is spent emulating
	// work. This is done to prevent the caller from noticing that the user was
	// not found (timing attack).
	user, err := h.UserService.GetUserByEmail(r.Context(), req.Username)
	if err == cargonaut.ErrUserNotFound {
		_, _ = password.Generate(h.secret, "fakework_invalid", password.DefaultCost)
		h.renderError(w, r, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Cancel the request if the user account is not active or the password hash
	// is not present in the user resource. If it is present, verify the
	// provided password.
	if user.Password == "" {
		h.renderError(w, r, http.StatusUnauthorized, err)
		return
	} else if err = password.Compare(h.secret, req.Password, user.Password); err == password.ErrPasswordMismatch {
		h.renderError(w, r, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Create authentication token for the user and store the newly created
	// token in the storage.
	token, err := jwt.NewToken(h.secret, user)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if err = h.UserService.CreateToken(r.Context(), token); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := loginResponse{
		Token:       token.Token,
		TokenExpiry: token.ExpiresAt,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}
}
