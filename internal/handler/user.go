package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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

func (h *Handler) listUserRatings(w http.ResponseWriter, r *http.Request) {
	if userID, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if ratings, err := h.UserRepository.ListRatings(r.Context(), userID); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, ratings)
	}
}

func (h *Handler) createUserRating(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	userID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	if uuid.Equal(authUserID, userID) {
		h.renderErrorf(w, r, http.StatusBadRequest, "can not rate yourself")
		return
	}

	var rating cargonaut.Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Author is the user who sent the request. User is the user the rating will
	// be given to.
	rating.UserID = userID
	rating.AuthorID = authUserID

	if err := h.UserRepository.CreateRating(r.Context(), &rating); err == cargonaut.ErrRatingExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) listUserVehicles(w http.ResponseWriter, r *http.Request) {
	if userID, err := uuid.FromString(chi.URLParam(r, "user_id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if vehicles, err := h.UserRepository.ListVehicles(r.Context(), userID); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, vehicles)
	}
}

func (h *Handler) getUserVehicle(w http.ResponseWriter, r *http.Request) {
	if userID, err := uuid.FromString(chi.URLParam(r, "user_id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if vehicleID, err := uuid.FromString(chi.URLParam(r, "vehicle_id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if vehicle, err := h.UserRepository.GetVehicle(r.Context(), userID, vehicleID); err == cargonaut.ErrVehicleNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, vehicle)
	}
}

func (h *Handler) createUserVehicle(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	userID, err := uuid.FromString(chi.URLParam(r, "user_id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Make sure we can not create a vehicle for other users by making sure the
	// user ID in the URI and auth token match.
	if !uuid.Equal(authUserID, userID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not create vehicle for another user")
		return
	}

	var vehicle cargonaut.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	vehicle.UserID = userID
	if err := h.UserRepository.CreateVehicle(r.Context(), &vehicle); err == cargonaut.ErrVehicleExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) updateUserVehicle(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	userID, err := uuid.FromString(chi.URLParam(r, "user_id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Make sure we can not update a vehicle of another users by making sure the
	// user ID in the URI and auth token match.
	if !uuid.Equal(authUserID, userID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not update vehicle of another user")
		return
	}

	vehicleID, err := uuid.FromString(chi.URLParam(r, "vehicle_id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	var vehicle cargonaut.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	vehicle.ID = vehicleID
	vehicle.UserID = userID
	if err := h.UserRepository.UpdateVehicle(r.Context(), &vehicle); err == cargonaut.ErrVehicleExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) deleteUserVehicle(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	userID, err := uuid.FromString(chi.URLParam(r, "user_id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Make sure we can not delete a vehicle for another users by making sure
	// the user ID in the URI and auth token match.
	if !uuid.Equal(authUserID, userID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not delete vehicle for another user")
		return
	}

	if vehicleID, err := uuid.FromString(chi.URLParam(r, "vehicle_id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if err := h.UserRepository.DeleteVehicle(r.Context(), userID, vehicleID); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

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
