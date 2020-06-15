package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
)

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	if users, err := h.UserService.ListUsers(r.Context()); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, users)
	}
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if user, err := h.UserService.GetUser(r.Context(), id); err == cargonaut.ErrUserNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, user)
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user cargonaut.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.UserService.CreateUser(r.Context(), &user); err == cargonaut.ErrUserExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	var user cargonaut.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	user.ID = uuid.FromStringOrNil(chi.URLParam(r, "id"))
	if err := h.UserService.UpdateUser(r.Context(), &user); err == cargonaut.ErrUserExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if err := h.UserService.DeleteUser(r.Context(), id); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}
