package cart

import (
	"testing"

	"be-ecommerce/internal/customerror"
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddToCart(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	cartID := uuid.New()

	tests := []struct {
		name          string
		userIDStr     string
		req           AddToCartRequest
		mockSetup     func(repo *MockCartRepository)
		expectedError *customerror.AppError
	}{
		{
			name:      "Success Add New Item",
			userIDStr: userID.String(),
			req: AddToCartRequest{
				ProductID: productID.String(),
				Quantity:  2,
			},
			mockSetup: func(repo *MockCartRepository) {
				repo.On("GetOrCreateCart", userID).Return(&models.Cart{ID: cartID, UserID: userID}, nil)
				repo.On("GetCartItem", cartID, productID).Return(nil, nil) // not exist
				repo.On("CreateCartItem", mock.AnythingOfType("*models.CartItem")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Success Update Existing Item",
			userIDStr: userID.String(),
			req: AddToCartRequest{
				ProductID: productID.String(),
				Quantity:  1,
			},
			mockSetup: func(repo *MockCartRepository) {
				repo.On("GetOrCreateCart", userID).Return(&models.Cart{ID: cartID, UserID: userID}, nil)
				itemID := uuid.New()
				repo.On("GetCartItem", cartID, productID).Return(&models.CartItem{ID: itemID}, nil)
				repo.On("UpdateCartItemQuantity", itemID, 1).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCartRepository)
			tt.mockSetup(mockRepo)
			useCase := NewCartUseCase(mockRepo)

			err := useCase.AddToCart(tt.userIDStr, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
