package health

import (
	"net/http"

	"github.com/mzeahmed/gobooking/internal/response"
)

// Handler handles all HTTP requests related to the health module.
type Handler struct {
	service *Service
}

// NewHandler creates a new health handler.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Health returns the application health status.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {

	resp := h.service.Health()

	response.JSON(w, http.StatusOK, resp)
}
