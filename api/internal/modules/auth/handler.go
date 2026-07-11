package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mzeahmed/gobooking/internal/modules/user"
	"github.com/mzeahmed/gobooking/internal/response"
)

// Handler handles all HTTP requests related to the auth module.
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler.
func NewHandler(service *Service) *Handler {

	return &Handler{
		service: service,
	}
}

// Register handles POST /auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, user.ErrEmailTaken) {
			response.JSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
			return
		}

		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	response.JSON(w, http.StatusCreated, res)
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := req.Validate(); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
			return
		}

		response.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	response.JSON(w, http.StatusOK, res)
}
