package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
)

func (h *Handler) listTrips(w http.ResponseWriter, r *http.Request) {
	if trips, err := h.TripRepository.ListTrips(r.Context()); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, trips)
	}
}

func (h *Handler) getTrip(w http.ResponseWriter, r *http.Request) {
	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if trip, err := h.TripRepository.GetTrip(r.Context(), id); err == cargonaut.ErrTripNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, trip)
	}
}

func (h *Handler) createTrip(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	var trip cargonaut.Trip
	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	trip.UserID = authUserID
	if err := h.TripRepository.CreateTrip(r.Context(), &trip); err == cargonaut.ErrTripExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) updateTrip(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	var trip cargonaut.Trip
	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Make sure we can not update a trip of another users by making sure the
	// user ID of the stored trip and the authenticated user match.
	if trip, err := h.TripRepository.GetTrip(r.Context(), id); err == cargonaut.ErrTripNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if !uuid.Equal(trip.UserID, authUserID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not update trip of another user")
		return
	}

	trip.UserID = authUserID
	trip.ID = uuid.FromStringOrNil(chi.URLParam(r, "id"))
	if err := h.TripRepository.UpdateTrip(r.Context(), &trip); err == cargonaut.ErrTripExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) deleteTrip(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	// Make sure we can not delete a trip of other users by making sure the user
	// ID of the stored trip and the authenticated user match.
	if trip, err := h.TripRepository.GetTrip(r.Context(), id); err == cargonaut.ErrTripNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if !uuid.Equal(trip.UserID, authUserID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not delete trip of another user")
		return
	}

	if err := h.TripRepository.DeleteTrip(r.Context(), id); err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}

func (h *Handler) getTripRating(w http.ResponseWriter, r *http.Request) {
	if id, err := uuid.FromString(chi.URLParam(r, "id")); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
	} else if rating, err := h.TripRepository.GetTrip(r.Context(), id); err == cargonaut.ErrRatingNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		h.renderOK(w, r, rating)
	}
}

func (h *Handler) createTripRating(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := h.userIDFromRequest(r.Context(), w, r)
	if !ok {
		return
	}

	tripID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	var rating cargonaut.Rating
	if err = json.NewDecoder(r.Body).Decode(&rating); err != nil {
		h.renderError(w, r, http.StatusBadRequest, err)
		return
	}

	trip, err := h.TripRepository.GetTrip(r.Context(), tripID)
	if err == cargonaut.ErrTripNotFound {
		h.renderError(w, r, http.StatusNotFound, err)
		return
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return
	} else if trip.RiderID != nil && !uuid.Equal(*trip.RiderID, authUserID) {
		h.renderErrorf(w, r, http.StatusForbidden, "can not rate trip taken by another user")
		return
	}

	// Author is the user who sent the request. User is the user the rating will
	// be given to.
	rating.UserID = trip.UserID
	rating.AuthorID = authUserID
	rating.TripID = tripID

	if err := h.TripRepository.CreateRating(r.Context(), &rating); err == cargonaut.ErrRatingExists {
		h.renderError(w, r, http.StatusConflict, err)
	} else if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
	} else {
		render.NoContent(w, r)
	}
}
