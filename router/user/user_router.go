package userRouter

import (
	"GoPattern/handlers"
	"GoPattern/middleware"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handlers.UserHandler
}

func NewModule(handler *handlers.UserHandler) *Module {
	return &Module{handler: handler}
}

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth)

		r.Get("/", m.handler.GetUsers)
		r.Get("/{id}", m.handler.GetUserByID)
	})

	return r
}
