package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/krishnasingh0210/blog/internal/api/dto"
	"github.com/krishnasingh0210/blog/internal/domain"
	"github.com/krishnasingh0210/blog/pkg/response"
)

type AuthHandler struct {
	userService domain.UserService
	validate    *validator.Validate
}

func NewAuthHandler(userService domain.UserService) *AuthHandler {
	return &AuthHandler{userService: userService, validate: validator.New()}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.userService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateData) {
			response.Error(w, http.StatusConflict, "email already registered")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	response.JSON(w, http.StatusCreated, dto.UserResponse{
		ID: u.ID.String(), Name: u.Name, Email: u.Email, Role: string(u.Role),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCreds) {
			response.Error(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to login")
		return
	}

	response.JSON(w, http.StatusOK, dto.LoginResponse{AccessToken: token})
}