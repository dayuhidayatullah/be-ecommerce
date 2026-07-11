package product

import (
	"encoding/json"
	"net/http"

	"be-ecommerce/internal/utils"

	"github.com/go-playground/validator/v10"
)

type HttpHandler struct {
	useCase  ProductUseCase
	validate *validator.Validate
}

// NewHttpHandler creates a new instance of HttpHandler
func NewHttpHandler(useCase ProductUseCase, validate *validator.Validate) *HttpHandler {
	return &HttpHandler{
		useCase:  useCase,
		validate: validate,
	}
}

func (h *HttpHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	resp, err := h.useCase.CreateCategory(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusCreated, "Category created successfully", resp)
}

func (h *HttpHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	resp, err := h.useCase.GetCategories()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve categories", err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Categories retrieved successfully", resp)
}

func (h *HttpHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	resp, err := h.useCase.CreateProduct(req)
	if err != nil {
		if err.Error() == "category not found" {
			utils.RespondError(w, http.StatusNotFound, "Category not found", "not_found")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create product", err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusCreated, "Product created successfully", resp)
}

func (h *HttpHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	resp, err := h.useCase.GetProducts()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve products", err.Error())
		return
	}

	utils.RespondSuccess(w, http.StatusOK, "Products retrieved successfully", resp)
}
