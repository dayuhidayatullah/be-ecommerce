package product

import (
	"be-ecommerce/internal/models"
	"github.com/google/uuid"
)

// DTOs
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateProductRequest struct {
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

type ProductResponse struct {
	ID          string           `json:"id"`
	Category    CategoryResponse `json:"category"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
}

// Interfaces
type ProductRepository interface {
	CreateCategory(category *models.Category) error
	GetCategories() ([]models.Category, error)
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	
	CreateProduct(product *models.Product) error
	GetProducts() ([]models.Product, error)
}

type ProductUseCase interface {
	CreateCategory(req CreateCategoryRequest) (CategoryResponse, error)
	GetCategories() ([]CategoryResponse, error)
	
	CreateProduct(req CreateProductRequest) (ProductResponse, error)
	GetProducts() ([]ProductResponse, error)
}
