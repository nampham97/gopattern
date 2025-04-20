package MDRouter

import (
	"GoPattern/handlers"
	"GoPattern/middleware"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *handlers.MasterdataHandler
}

func NewModule(handler *handlers.MasterdataHandler) *Module {
	return &Module{handler: handler}
}

func (m *Module) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.JWTAuth)

	r.Post("/getCountries", m.handler.GetCountries)
	r.Post("/getProvinces", m.handler.GetProvinces)

	return r
}
