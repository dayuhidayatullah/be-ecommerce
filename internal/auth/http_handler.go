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

	// 1. Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// 2. Validate input
	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// 3. Process via UseCase
	resp, err := h.useCase.Register(req)
	if err != nil {
		if err.Error() == "email_conflict" {
			utils.RespondError(w, http.StatusConflict, "Email already registered", "email_conflict")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Registration failed", err.Error())
		return
	}

	// 4. Return response
	utils.RespondSuccess(w, http.StatusCreated, "User registered successfully", resp)
}

// Login authenticates a user and returns a JWT token
func (h *HttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// 1. Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	// 2. Validate input
	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// 3. Process via UseCase
	resp, err := h.useCase.Login(req)
	if err != nil {
		if err.Error() == "invalid_credentials" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid email or password", "invalid_credentials")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Login failed", err.Error())
		return
	}

	// 4. Return response
	utils.RespondSuccess(w, http.StatusOK, "Login successful", resp)
}
