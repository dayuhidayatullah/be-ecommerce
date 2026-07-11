package cart

import (
	"encoding/json"
	"net/http"

	"be-ecommerce/internal/middleware"
	"be-ecommerce/internal/utils"

	"github.com/go-playground/validator/v10"
)

type HttpHandler struct {
	useCase  CartUseCase
	validate *validator.Validate
}

// NewHttpHandler creates a new instance of HttpHandler
func NewHttpHandler(useCase CartUseCase, validate *validator.Validate) *HttpHandler {
	return &HttpHandler{
		useCase:  useCase,
		validate: validate,
	}
}

func (h *HttpHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// Get UserID from Context (set by RequireAuth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized", "missing user context")
		return
	}

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	if err := h.useCase.AddToCart(userID, req); err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Item added to cart successfully", nil)
}

func (h *HttpHandler) GetMyCart(w http.ResponseWriter, r *http.Request) {
	// Get UserID from Context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized", "missing user context")
		return
	}

	resp, err := h.useCase.GetMyCart(userID)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Cart retrieved successfully", resp)
}
