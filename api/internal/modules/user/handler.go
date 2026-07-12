package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mzeahmed/gobooking/internal/response"
)

// Handler handles all HTTP requests related to the user module.
type Handler struct {
	service *Service
}

// NewHandler creates a new user handler.
func NewHandler(service *Service) *Handler {

	return &Handler{
		service: service,
	}
}

// Delete handles DELETE /users.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	var req DeleteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})

		return
	}

	if err := req.Validate(); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})

		return
	}

	if err := h.service.DeleteUser(r.Context(), req); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.JSON(w, http.StatusNotFound, map[string]string{
				"error": "user not found",
			})

			return
		}

		response.JSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})

		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
