package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/rakyll/statik/fs"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
	_ "github.com/my-cargonaut/cargonaut/internal/ui" // UI
	"github.com/my-cargonaut/cargonaut/pkg/version"
)

var _ http.Handler = (*Handler)(nil)

// Handler provides all hhtp handlers.
type Handler struct {
	log    *log.Logger
	router chi.Router

	secret []byte

	TripRepository    cargonaut.TripRepository
	UserRepository    cargonaut.UserRepository
	VehicleRepository cargonaut.VehicleRepository
	TokenBlacklist    cargonaut.TokenBlacklist
}

// NewHandler creates a new set of handlers.
func NewHandler(log *log.Logger, secret []byte) (*Handler, error) {
	h := &Handler{
		log:    log,
		router: chi.NewRouter(),

		secret: secret,
	}

	// Default NotFound and MethodNotAllowed handlers. The NotFound handler
	// redirects to the index page.
	h.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.RedirectHandler("/", http.StatusPermanentRedirect).ServeHTTP(w, r)
	})
	h.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
	})

	// CORS setup.
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	})

	// Base middleware stack: CORS, compression, panic recoverer, request
	// timeout.
	h.router.Use(cors.Handler)
	h.router.Use(middleware.Compress(5))
	h.router.Use(middleware.Recoverer)
	h.router.Use(middleware.Timeout(30 * time.Second))

	// Serve user interface.
	ui, err := fs.NewWithNamespace("ui")
	if err != nil {
		return nil, fmt.Errorf("create web ui file system: %w", err)
	}
	h.router.Handle("/*", http.FileServer(ui))

	// Serve API.
	h.router.Route("/api/v1", func(api chi.Router) {
		// NotFound and MethodNotAllowed handlers for API.
		api.NotFound(func(w http.ResponseWriter, r *http.Request) {
			code := http.StatusNotFound
			h.renderError(w, r, code, errors.New(http.StatusText(code)))
		})
		api.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			code := http.StatusMethodNotAllowed
			h.renderError(w, r, code, errors.New(http.StatusText(code)))
		})

		// API middleware (JSON content type & renderer).
		api.Use(middleware.AllowContentType("application/json"))
		api.Use(render.SetContentType(render.ContentTypeJSON))

		// Authentication routes.
		api.Post("/auth/login", h.login)
		api.Patch("/auth/refresh", h.refresh)
		api.Post("/auth/logout", h.logout)
		api.Post("/auth/register", h.register)

		// Special user profile picture route.
		api.Get("/users/{id}/avatar", h.getUserAvatar)

		// Authenticated routes.
		api.Group(func(r chi.Router) {
			// Authentication.
			r.Use(jwtauth.Verifier(jwtauth.New("HS256", secret, nil)))
			r.Use(jwtauth.Authenticator)

			// Trip API.
			r.Get("/trips", h.listTrips)
			r.Get("/trips/{id}", h.getTrip)
			r.Post("/trips", h.createTrip)
			r.Put("/trips/{id}", h.updateTrip)
			r.Delete("/trips/{id}", h.deleteTrip)
			r.Get("/trips/{id}/ratings", h.getTripRating)
			r.Post("/trips/{id}/ratings", h.createTripRating)

			// User API.
			// r.Get("/users", h.listUsers)
			r.Get("/users/{id}", h.getUser)
			// r.Post("/users", h.createUser)
			// r.Put("/users/{id}", h.updateUser)
			// r.Delete("/users/{id}", h.deleteUser)
			r.Get("/users/{id}/ratings", h.listUserRatings)
			r.Get("/users/{id}/vehicles", h.listUserVehicles)
			r.Post("/users/{user_id}/trips/{trip_id}", h.bookTrip)
			r.Put("/users/{user_id}/trips/{trip_id}", h.cancelTrip)

			// Vehicle API.
			r.Get("/vehicles", h.listVehicles)
			r.Get("/vehicles/{id}", h.getVehicle)
			r.Post("/vehicles", h.createVehicle)
			r.Put("/vehicles/{id}", h.updateVehicle)
			r.Delete("/vehicles/{id}", h.deleteVehicle)
		})
	})

	return h, nil
}

// ServeHTTP serves all http routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", fmt.Sprintf("cargonaut %s", version.Release()))
	h.router.ServeHTTP(w, r)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	render.Status(r, code)
	render.JSON(w, r, data)
}

func (h *Handler) renderOK(w http.ResponseWriter, r *http.Request, data interface{}) {
	h.render(w, r, http.StatusOK, data)
}

func (h *Handler) renderError(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.log.Printf("[%s %s]: %s", r.Method, r.RequestURI, err)

	h.render(w, r, code, map[string]string{
		"status": http.StatusText(code),
		"error":  err.Error(),
	})
}

func (h *Handler) renderErrorf(w http.ResponseWriter, r *http.Request, code int, format string, a ...interface{}) {
	err := fmt.Errorf(format, a...)
	h.renderError(w, r, code, err)
}

func (h *Handler) userIDFromRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return uuid.Nil, false
	}

	userClaim, ok := claims["user"]
	if !ok {
		h.renderErrorf(w, r, http.StatusInternalServerError, "user claim missing")
		return uuid.Nil, false
	}

	user, ok := userClaim.(map[string]interface{})
	if !ok {
		h.renderErrorf(w, r, http.StatusInternalServerError, "could not get user claim missing")
		return uuid.Nil, false
	}

	if _, ok = user["id"]; !ok {
		h.renderErrorf(w, r, http.StatusInternalServerError, "user.id claim missing")
		return uuid.Nil, false
	}

	userIDClaim, ok := user["id"].(string)
	if !ok {
		h.renderErrorf(w, r, http.StatusInternalServerError, "user.id claim is not a string")
		return uuid.Nil, false
	}

	userID, err := uuid.FromString(userIDClaim)
	if err != nil {
		h.renderError(w, r, http.StatusInternalServerError, err)
		return uuid.Nil, false
	}

	return userID, true
}
