package product

import (
	"errors"

	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productUseCase struct {
	repo ProductRepository
}

// NewProductUseCase creates a new instance of ProductUseCase
func NewProductUseCase(repo ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (u *productUseCase) CreateCategory(req CreateCategoryRequest) (CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := u.repo.CreateCategory(&category); err != nil {
		return CategoryResponse{}, err
	}

	return CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (u *productUseCase) GetCategories() ([]CategoryResponse, error) {
	categories, err := u.repo.GetCategories()
	if err != nil {
		return nil, err
	}

	var responses []CategoryResponse
	for _, c := range categories {
		responses = append(responses, CategoryResponse{
			ID:          c.ID.String(),
			Name:        c.Name,
			Description: c.Description,
		})
	}
	return responses, nil
}

func (u *productUseCase) CreateProduct(req CreateProductRequest) (ProductResponse, error) {
	// Parse CategoryID to UUID
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return ProductResponse{}, errors.New("invalid category id format")
	}

	// Verify category exists
	category, err := u.repo.GetCategoryByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProductResponse{}, errors.New("category not found")
		}
		return ProductResponse{}, err
	}

	product := models.Product{
		CategoryID:  categoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := u.repo.CreateProduct(&product); err != nil {
		return ProductResponse{}, err
	}

	return ProductResponse{
		ID:          product.ID.String(),
		Category: CategoryResponse{
			ID:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
		},
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}, nil
}

func (u *productUseCase) GetProducts() ([]ProductResponse, error) {
	products, err := u.repo.GetProducts()
	if err != nil {
		return nil, err
	}

	var responses []ProductResponse
	for _, p := range products {
		responses = append(responses, ProductResponse{
			ID: p.ID.String(),
			Category: CategoryResponse{
				ID:          p.Category.ID.String(),
				Name:        p.Category.Name,
				Description: p.Category.Description,
			},
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
		})
	}
	return responses, nil
}
