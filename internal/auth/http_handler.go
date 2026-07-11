package auth

import (
	"encoding/json"
	"net/http"

	"be-ecommerce/internal/utils"

	"github.com/go-playground/validator/v10"
)

// HttpHandler handles HTTP requests for authentication
type HttpHandler struct {
	useCase  AuthUseCase
	validate *validator.Validate
}

// NewHttpHandler creates a new instance of HttpHandler
func NewHttpHandler(useCase AuthUseCase, validate *validator.Validate) *HttpHandler {
	return &HttpHandler{
		useCase:  useCase,
		validate: validate,
	}
}

// Register creates a new user account
func (h *HttpHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	resp, err := h.useCase.Register(req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.RespondSuccess(w, http.StatusCreated, "User registered successfully", resp)
}

// Login authenticates a user and returns a JWT token
func (h *HttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	resp, err := h.useCase.Login(req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Login successful", resp)
}
