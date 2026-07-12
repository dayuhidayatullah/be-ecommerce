package product

import (
	"testing"

	"be-ecommerce/internal/customerror"
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateCategory(t *testing.T) {
	tests := []struct {
		name          string
		req           CreateCategoryRequest
		mockSetup     func(repo *MockProductRepository)
		expectedError *customerror.AppError
	}{
		{
			name: "Success Create Category",
			req: CreateCategoryRequest{
				Name:        "Electronics",
				Description: "Electronic items",
			},
			mockSetup: func(repo *MockProductRepository) {
				repo.On("CreateCategory", mock.AnythingOfType("*models.Category")).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			tt.mockSetup(mockRepo)
			useCase := NewProductUseCase(mockRepo)

			resp, err := useCase.CreateCategory(tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Name, resp.Name)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateProduct(t *testing.T) {
	categoryID := uuid.New()
	tests := []struct {
		name          string
		req           CreateProductRequest
		mockSetup     func(repo *MockProductRepository)
		expectedError *customerror.AppError
	}{
		{
			name: "Success Create Product",
			req: CreateProductRequest{
				CategoryID:  categoryID.String(),
				Name:        "Laptop",
				Description: "Gaming Laptop",
				Price:       1500,
				Stock:       10,
				ImageURL:    "http://example.com/laptop.png",
			},
			mockSetup: func(repo *MockProductRepository) {
				repo.On("GetCategoryByID", categoryID).Return(&models.Category{ID: categoryID, Name: "Electronics"}, nil)
				repo.On("CreateProduct", mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed - Category Not Found",
			req: CreateProductRequest{
				CategoryID:  categoryID.String(),
				Name:        "Laptop",
				Price:       1500,
				Stock:       10,
			},
			mockSetup: func(repo *MockProductRepository) {
				repo.On("GetCategoryByID", categoryID).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: customerror.NewNotFoundError("category not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			tt.mockSetup(mockRepo)
			useCase := NewProductUseCase(mockRepo)

			resp, err := useCase.CreateProduct(tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				appErr, _ := err.(*customerror.AppError)
				assert.Equal(t, tt.expectedError.Type, appErr.Type)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Name, resp.Name)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
