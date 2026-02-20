package handlers

import (
	"errors"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/dto"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterReq true "User registration data"
// @Success      201 {object} helper.Response{data=dto.AuthResp}
// @Failure      400 {object} helper.Response "Invalid request data or user already exists"
// @Router       /auth/signup [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRegisterReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid given payload")
		return
	}

	user, err := h.authService.Register(r.Context(), &payload)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.BadRequestResponse(w, "Registration failed", err)
		case errors.Is(err, repository.ErrDuplicateEmail):
			helper.EditConflictResponse(w, "Registration failed", err)
		default:
			helper.InternalServerError(w, "Failed to register user", err)
		}
		return
	}

	helper.CreatedResponse(w, "User successfully registered", user)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginReq true "User login credentials"
// @Success      200 {object} helper.Response{data=dto.AuthResp} "Login successfully"
// @Failure      401 {object} helper.Response "Invalid credentials"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateLoginReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid given payload")
		return
	}

	login, err := h.authService.Login(r.Context(), &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to login", err)
		return
	}

	helper.SuccessResponse(w, "User successfully logged in", login)
}

// RefreshToken docs
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenReq true "Refresh token"
// @Success 200 {object} helper.Response{data=dto.AuthResp} "Token refreshed successfully"
// @Failure 401 {object} helper.Response "Invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload dto.RefreshTokenReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRefreshTokenReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid given payload")
		return
	}

	refreshToken, err := h.authService.RefreshToken(r.Context(), &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to refresh token", err)
		return
	}

	helper.SuccessResponse(w, "Refresh token successfully generated", refreshToken)
}

// Logout docs
// @Summary User logout
// @Description Invalidate refresh token and logout user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenReq true "Refresh token to invalidate"
// @Success 200 {object} helper.Response "Logout successful"
// @Failure 400 {object} helper.Response "Invalid request data"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var payload dto.RefreshTokenReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRefreshTokenReq(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid given payload")
		return
	}

	if err := h.authService.Logout(r.Context(), &payload); err != nil {
		helper.InternalServerError(w, "Failed to logout", err)
		return
	}

	helper.SuccessResponse(w, "Successfully logged out", nil)
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
