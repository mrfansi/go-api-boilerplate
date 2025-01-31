package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mrfansi/go-api-boilerplate/internal/application/service"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/errors"
)

type UserHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user"`
}

// CreateUser godoc
// @Summary Create new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation request"
// @Success 201 {object} entity.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrValidation)
		return
	}

	user, err := h.userService.Create(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		switch err {
		case errors.ErrUserAlreadyExists:
			respondWithError(w, http.StatusConflict, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entity.User
// @Failure 404 {object} errors.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User update request"
// @Success 200 {object} entity.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrValidation)
		return
	}

	user, err := h.userService.Update(r.Context(), id, req.Name)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} errors.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if err := h.userService.Delete(r.Context(), id); err != nil {
		switch err {
		case errors.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUsers godoc
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} entity.User
// @Router /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	users, total, err := h.userService.List(r.Context(), page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		return
	}

	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	respondWithJSON(w, http.StatusOK, users)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user's password
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body ChangePasswordRequest true "Password change request"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /users/{id}/password [put]
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrValidation)
		return
	}

	if err := h.userService.ChangePassword(r.Context(), id, req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case errors.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err)
		case errors.ErrInvalidPassword:
			respondWithError(w, http.StatusBadRequest, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateRole godoc
// @Summary Update user role
// @Description Update user's role
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateRoleRequest true "Role update request"
// @Success 204 "No Content"
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /users/{id}/role [put]
func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrValidation)
		return
	}

	if err := h.userService.UpdateRole(r.Context(), id, req.Role); err != nil {
		switch err {
		case errors.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err)
		case errors.ErrInvalidRole:
			respondWithError(w, http.StatusBadRequest, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
