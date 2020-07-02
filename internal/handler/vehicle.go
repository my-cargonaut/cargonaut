package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
)

// func (h *Handler) listVehicles(w http.ResponseWriter, r *http.Request) {
// 	if vehicles, err := h.VehicleRepository.ListVehicles(r.Context()); err != nil {
// 		h.renderError(w, r, http.StatusInternalServerError, err)
// 	} else {
// 		h.renderOK(w, r, vehicles)
// 	}
// }

func (h *Handler) getVehicle(w http.ResponseWriter, r *http.Request) {
	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if vehicle, err := h.VehicleRepository.GetVehicle(r.Context(), id); err == cargonaut.ErrVehicleNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, vehicle)
	}
}

func (h *Handler) createVehicle(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	var vehicle cargonaut.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	vehicle.UserID = authUserID
	if err := h.VehicleRepository.CreateVehicle(r.Context(), &vehicle); err == cargonaut.ErrVehicleExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) updateVehicle(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	var vehicle cargonaut.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Make sure we can not update a vehicle of another users by making sure the
	// user ID of the stored vehicle and the authenticated user match.
	if vehicle, err := h.VehicleRepository.GetVehicle(r.Context(), id); err == cargonaut.ErrVehicleNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if !uuid.Equal(vehicle.UserID, authUserID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not update vehicle of another user")
		return
	}

	vehicle.ID = id
	vehicle.UserID = authUserID
	if err := h.VehicleRepository.UpdateVehicle(r.Context(), &vehicle); err == cargonaut.ErrVehicleExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) deleteVehicle(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Make sure we can not delete a vehicle of other users by making sure the user
	// ID of the stored vehicle and the authenticated user match.
	if vehicle, err := h.VehicleRepository.GetVehicle(r.Context(), id); err == cargonaut.ErrVehicleNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if !uuid.Equal(vehicle.UserID, authUserID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not delete vehicle of another user")
		return
	}

	if err := h.VehicleRepository.DeleteVehicle(r.Context(), id); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}
