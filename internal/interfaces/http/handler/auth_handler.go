package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mrfansi/go-api-boilerplate/internal/application/service"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/errors"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, errors.ErrValidation)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case errors.ErrInvalidCredential:
			respondWithError(w, http.StatusUnauthorized, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, TokenResponse{Token: token})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Get new JWT token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} TokenResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, errors.ErrInvalidToken)
		return
	}

	newToken, err := h.authService.RefreshToken(token)
	if err != nil {
		switch err {
		case errors.ErrInvalidToken, errors.ErrTokenExpired:
			respondWithError(w, http.StatusUnauthorized, err)
		default:
			respondWithError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, TokenResponse{Token: newToken})
}

// Helper functions

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJSON(w, code, errors.NewErrorResponse(code, err.Error()))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
