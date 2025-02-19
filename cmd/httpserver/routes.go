package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"orderservice/state"
	"time"
)

func routes(s *state.State) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(120 * time.Second))

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsOptions).Handler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", HandleRegisterUser(s))
		r.Post("/login", HandleLogin(s))
		r.Post("/token/refresh", HandleRefreshToken(s))
		r.With(AuthMiddleware(s)).Post("/logout", HandleLogout(s))
	})

	r.Route("/api/v1/orders", func(r chi.Router) {
		r.Use(AuthMiddleware(s))
		r.Post("/", HandleCreateOrder(s))
		r.Get("/all", HandlerGetAllOrders(s))
		r.Put("/{id}/cancel", HandleCancelOrder(s))
	})

	return r
}
