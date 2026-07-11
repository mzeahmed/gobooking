package health

import (
	"net/http"

	"github.com/mzeahmed/gobooking/internal/middleware"
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

// Protected returns the application health status along with the
// identity of the authenticated caller. It exists as a usage example of
// the Authenticate middleware, guarded by it in RegisterRoutes.
func (h *Handler) Protected(w http.ResponseWriter, r *http.Request) {

	authUser, _ := middleware.UserFromContext(r.Context())

	response.JSON(w, http.StatusOK, ProtectedResponse{
		Status: "ok",
		UserID: authUser.ID,
		Roles:  authUser.Roles,
	})
}
