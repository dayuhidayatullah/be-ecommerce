package order

import (
	"encoding/json"
	"net/http"

	"be-ecommerce/internal/middleware"
	"be-ecommerce/internal/utils"

	"github.com/go-playground/validator/v10"
)

type HttpHandler struct {
	useCase  OrderUseCase
	validate *validator.Validate
}

// NewHttpHandler creates a new instance of HttpHandler
func NewHttpHandler(useCase OrderUseCase, validate *validator.Validate) *HttpHandler {
	return &HttpHandler{
		useCase:  useCase,
		validate: validate,
	}
}

func (h *HttpHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	// Get UserID from Context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized", "missing user context")
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	resp, err := h.useCase.Checkout(userID, req)
	if err != nil {
		if err.Error() == "cart is empty" {
			utils.RespondError(w, http.StatusBadRequest, "Cannot checkout an empty cart", "empty_cart")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Checkout failed", err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusCreated, "Checkout successful, order created", resp)
}

func (h *HttpHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	// Get UserID from Context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized", "missing user context")
		return
	}

	resp, err := h.useCase.GetMyOrders(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve orders", err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Orders retrieved successfully", resp)
}
