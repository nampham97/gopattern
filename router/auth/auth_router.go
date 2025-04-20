package authRouter

import (
	"GoPattern/handlers"
	"GoPattern/middleware"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handlers.AuthHandler
}

func NewModule(handler *handlers.AuthHandler) *Module {
	return &Module{handler: handler}
}

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", m.handler.Login)
	r.Post("/refresh", m.handler.RefreshToken)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)
		r.Post("/logout", m.handler.Logout)
	})

	return r
}
