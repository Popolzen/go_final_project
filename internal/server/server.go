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

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)

		r.Route("/secrets", func(r chi.Router) {
			r.Use(h.AuthMiddleware)

			r.Post("/", h.CreateSecret)
			r.Get("/", h.ListSecrets)
			r.Get("/{id}", h.GetSecret)
			r.Put("/{id}", h.UpdateSecret)
			r.Delete("/{id}", h.DeleteSecret)
			r.Get("/by-name/{name}", h.GetSecretByName)
			r.Get("/changes", h.GetSecretsAfter)
		})
	})

	return r
}
