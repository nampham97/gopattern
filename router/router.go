package router

import (
	"GoPattern/db"
	"GoPattern/handlers"
	"GoPattern/middleware"
	"GoPattern/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	userRepo := repository.NewUserRepository(db.GetDB())
	userHandler := handlers.NewUserHandler(userRepo)

	r.Post("/login", handlers.Login)
	r.Post("/refresh", handlers.RefreshToken)

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUserByID)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Get("/", handlers.AdminOnly)
	})

	return r
}
