package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Routes настраивает все маршруты
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Post("/api/v1/auth/register", h.Register)
	r.Post("/api/v1/auth/login", h.Login)

	r.Post("/api/v1/secrets", h.CreateSecret)
	r.Get("/api/v1/secrets", h.ListSecrets)

	r.Get("/api/v1/secrets/{id}", h.GetSecret)
	r.Put("/api/v1/secrets/{id}", h.UpdateSecret)
	r.Delete("/api/v1/secrets/{id}", h.DeleteSecret)

	r.Get("/api/v1/secrets/by-name/{name}", h.GetSecretByName)

	r.Get("/api/v1/secrets/changes", h.GetSecretsAfter)

	return r
}
