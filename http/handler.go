package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rakyll/statik/fs"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
	_ "github.com/my-cargonaut/cargonaut/http/ui" // UI
	"github.com/my-cargonaut/cargonaut/version"
)

var _ http.Handler = (*Handler)(nil)

// Handler provides all hhtp handlers.
type Handler struct {
	log    *log.Logger
	router chi.Router

	UserService cargonaut.UserService
}

// NewHandler creates a new set of handlers.
func NewHandler(log *log.Logger) (*Handler, error) {
	h := &Handler{
		log:    log,
		router: chi.NewRouter(),
	}

	h.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.RedirectHandler("/", http.StatusPermanentRedirect).ServeHTTP(w, r)
	})
	h.router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
	})

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	})

	h.router.Use(cors.Handler)
	h.router.Use(middleware.Compress(5))
	h.router.Use(middleware.Recoverer)

	ui, err := fs.NewWithNamespace("ui")
	if err != nil {
		return nil, fmt.Errorf("create web ui file system: %w", err)
	}
	h.router.Handle("/*", http.FileServer(ui))

	h.router.Route("/api/v1", func(api chi.Router) {
		api.NotFound(func(w http.ResponseWriter, r *http.Request) {
			code := http.StatusNotFound
			h.renderError(w, r, code, errors.New(http.StatusText(code)))
		})
		api.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			code := http.StatusMethodNotAllowed
			h.renderError(w, r, code, errors.New(http.StatusText(code)))
		})

		api.Use(middleware.AllowContentType("application/json"))
		api.Use(render.SetContentType(render.ContentTypeJSON))

		api.Get("/users", h.listUsers)
		api.Get("/users/{id}", h.getUser)
		api.Post("/users", h.createUser)
		api.Put("/users/{id}", h.updateUser)
		api.Delete("/users/{id}", h.deleteUser)
	})

	return h, nil
}

// ServeHTTP serves all http routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", fmt.Sprintf("cargonaut %s", version.Release()))
	h.router.ServeHTTP(w, r)
}

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
	h.renderError(w, r, code, fmt.Errorf(format, a...))
}
