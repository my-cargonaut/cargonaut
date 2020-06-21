package handler

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
)

// func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
// 	if users, err := h.UserRepository.ListUsers(r.Context()); err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 	} else {
// 		h.renderOK(w, r, users)
// 	}
// }

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if user, err := h.UserRepository.GetUser(r.Context(), id); err == cargonaut.ErrUserNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, user)
	}
}

// func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
// 	var user cargonaut.User
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		h.renderError(w, r, http.StatusBadRequest, err)
// 		return
// 	}

// 	// Encrypt password.
// 	var err error
// 	if user.Password, err = password.Generate(h.secret, user.Password, password.DefaultCost); err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 		return
// 	}

// 	// Base64 decode avatar request data.
// 	imgSrc, err := base64.StdEncoding.DecodeString(user.Avatar)
// 	if err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 		return
// 	}

// 	// Decode image data.
// 	buf := bytes.NewBuffer(imgSrc)
// 	img, _, err := image.Decode(buf)
// 	if err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 		return
// 	}

// 	// Resize to a reasonable size.
// 	img = resize.Resize(250, 0, img, resize.Lanczos3)

// 	// Encode resized image as PNG and base64 encode.
// 	buf.Reset()
// 	if err := png.Encode(buf, img); err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 		return
// 	}
// 	user.Avatar = base64.StdEncoding.EncodeToString(buf.Bytes())

// 	if err := h.UserRepository.CreateUser(r.Context(), &user); err == cargonaut.ErrUserExists {
// 		h.renderError(w, r, http.StatusConflict, err)
// 	} else if err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 	} else {
// 		render.NoContent(w, r)
// 	}
// }

// func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
// 	var user cargonaut.User
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		h.renderError(w, r, http.StatusBadRequest, err)
// 		return
// 	}

// 	user.ID = uuid.FromStringOrNil(chi.URLParam(r, "id"))
// 	if err := h.UserRepository.UpdateUser(r.Context(), &user); err == cargonaut.ErrUserExists {
// 		h.renderError(w, r, http.StatusConflict, err)
// 	} else if err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 	} else {
// 		render.NoContent(w, r)
// 	}
// }

// func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
// 	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
// 		h.renderError(w, r, http.StatusBadRequest, err)
// 	} else if err := h.UserRepository.DeleteUser(r.Context(), id); err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 	} else {
// 		render.NoContent(w, r)
// 	}
// }

func (h *Handler) getUserAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")

	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if user, err := h.UserRepository.GetUser(r.Context(), id); err == cargonaut.ErrUserNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		if imgSrc, err := base64.StdEncoding.DecodeString(user.Avatar); err != nil {
			h.renderError(w, r, http.StatusInternalServerError, err)
		} else if _, err = io.Copy(w, bytes.NewReader(imgSrc)); err != nil {
			h.renderError(w, r, http.StatusInternalServerError, err)
		}
	}
}
