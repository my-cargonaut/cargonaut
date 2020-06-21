package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	_ "image/jpeg" // JPEG format
	"image/png"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/nfnt/resize"

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

type refreshRequest struct {
	Token string `json:"token"`
}

type refreshResponse struct {
	Token       string    `json:"token"`
	TokenExpiry time.Time `json:"token_expiry"`
}

type logoutRequest struct {
	Token string `json:"token"`
}

type registerRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	DisplayName string    `json:"display_name"`
	Birthday    time.Time `json:"birthday"`
	Avatar      string    `json:"avatar"`
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
	user, err := h.UserRepository.GetUserByEmail(r.Context(), req.Username)
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

	// Create an authentication token for the user and store the newly created
	// token in the storage.
	token, err := jwt.NewToken(h.secret, user)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if err = h.UserRepository.CreateToken(r.Context(), token); err != nil {
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

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	oldToken := &cargonaut.Token{
		Token: req.Token,
	}

	user, err := jwt.UserFromToken(h.secret, oldToken)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Check if the provided authentication token is blacklisted. If so, the
	// token is invalid and the request thus unauthenticated.
	// The old token is put on the token blacklist.
	var isBlacklisted bool
	if isBlacklisted, err = h.TokenBlacklist.IsTokenBlacklisted(r.Context(), oldToken.ID); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if isBlacklisted {
		h.renderError(w, r, http.StatusUnauthorized, err)
		return
	} else if err = h.TokenBlacklist.BlacklistToken(r.Context(), oldToken); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Delete the old token from the users token storage.
	if err = h.UserRepository.DeleteToken(r.Context(), user.ID, oldToken.ID); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Create a new authentication token for the user and store the newly
	// created token in the storage.
	token, err := jwt.NewToken(h.secret, user)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if err = h.UserRepository.CreateToken(r.Context(), token); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	resp := refreshResponse{
		Token:       token.Token,
		TokenExpiry: token.ExpiresAt,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	token := &cargonaut.Token{
		Token: req.Token,
	}

	user, err := jwt.UserFromToken(h.secret, token)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Check if the provided authentication token is blacklisted. If so, the
	// token is invalid and the request thus unauthenticated.
	// The token is put on the token blacklist.
	if isBlacklisted, err := h.TokenBlacklist.IsTokenBlacklisted(r.Context(), token.ID); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if isBlacklisted {
		h.renderError(w, r, http.StatusUnauthorized, err)
		return
	} else if err := h.TokenBlacklist.BlacklistToken(r.Context(), token); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Delete the token from the users token storage.
	if err := h.UserRepository.DeleteToken(r.Context(), user.ID, token.ID); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	h.renderOK(w, r, nil)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Encrypt password.
	var err error
	if req.Password, err = password.Generate(h.secret, req.Password, password.DefaultCost); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Base64 decode avatar request data.
	imgSrc, err := base64.StdEncoding.DecodeString(req.Avatar)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Decode image data.
	buf := bytes.NewBuffer(imgSrc)
	img, _, err := image.Decode(buf)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	// Resize to a reasonable size.
	img = resize.Resize(250, 0, img, resize.Lanczos3)

	// Encode resized image as PNG and base64 encode.
	buf.Reset()
	if err := png.Encode(buf, img); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}
	req.Avatar = base64.StdEncoding.EncodeToString(buf.Bytes())

	user := &cargonaut.User{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		Birthday:    req.Birthday,
		Avatar:      req.Avatar,
	}

	if err := h.UserRepository.CreateUser(r.Context(), user); err == cargonaut.ErrUserExists {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	}

	render.NoContent(w, r)
}
