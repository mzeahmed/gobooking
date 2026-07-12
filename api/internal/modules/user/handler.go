package user

import (
	"errors"
	"net/http"
	"slices"

	"github.com/mzeahmed/gobooking/internal/reqctx"
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

// List handles GET /users.
//
// This route is guarded by middleware.Authenticate (see
// Module.RegisterRoutes) and further restricted to callers with the
// "admin" role, since it exposes every registered user's email and roles.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	authUser, ok := reqctx.AuthUserFromContext(r.Context())

	if !ok {
		// Should not happen behind Authenticate, but fail closed if it does.
		response.JSON(w, http.StatusUnauthorized, map[string]string{
			"error": "authentication required",
		})

		return
	}

	if !slices.Contains(authUser.Roles, string(RoleAdmin)) {
		response.JSON(w, http.StatusForbidden, map[string]string{
			"error": "admin role required",
		})

		return
	}

	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error : " + err.Error(),
		})

		return
	}

	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = newUserResponse(u)
	}

	response.JSON(w, http.StatusOK, res)
}

// Delete handles DELETE /users/delete.
//
// This route is guarded by middleware.Authenticate (see
// Module.RegisterRoutes), so it only ever runs with a valid, authenticated
// caller. The target user ID is taken from that authenticated identity,
// not from the request body, so callers can only ever delete their own
// account.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := reqctx.AuthUserFromContext(r.Context())

	if !ok {
		// Should not happen behind Authenticate, but fail closed if it does.
		response.JSON(w, http.StatusUnauthorized, map[string]string{
			"error": "authentication required",
		})

		return
	}

	req := DeleteRequest{ID: authUser.ID}

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
			"error": "internal server error : " + err.Error(),
		})

		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}
