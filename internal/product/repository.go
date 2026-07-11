package product

import (
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *productRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *productRepository) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ?", id).Error
	return &category, err
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetProducts() ([]models.Product, error) {
	var products []models.Product
	// Preload the Category data so it's included in the product response
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}
