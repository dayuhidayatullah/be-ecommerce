package product

import (
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateCategory(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockProductRepository) GetCategories() ([]models.Category, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetProducts() ([]models.Product, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}
