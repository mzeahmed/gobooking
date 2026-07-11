package health

import "net/http"

type Module struct {
	handler *Handler
}

func New() *Module {

	service := NewService()

	handler := NewHandler(service)

	return &Module{
		handler: handler,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /health", m.handler.Health)
}
