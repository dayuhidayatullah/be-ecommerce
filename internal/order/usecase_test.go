package order

import (
	"errors"
	"testing"

	"be-ecommerce/internal/customerror"
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCheckout(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	tests := []struct {
		name          string
		userIDStr     string
		req           CheckoutRequest
		mockSetup     func(repo *MockOrderRepository)
		expectedError *customerror.AppError
	}{
		{
			name:      "Success Checkout",
			userIDStr: userID.String(),
			req: CheckoutRequest{
				ShippingAddress: "Jakarta",
			},
			mockSetup: func(repo *MockOrderRepository) {
				repo.On("Checkout", userID, "Jakarta").Return(&models.Order{
					ID:              orderID,
					UserID:          userID,
					ShippingAddress: "Jakarta",
					TotalAmount:     1000,
					Status:          "pending",
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:      "Failed - Empty Cart",
			userIDStr: userID.String(),
			req: CheckoutRequest{
				ShippingAddress: "Jakarta",
			},
			mockSetup: func(repo *MockOrderRepository) {
				repo.On("Checkout", userID, "Jakarta").Return(nil, errors.New("cart is empty"))
			},
			expectedError: customerror.NewBadRequestError("cannot checkout an empty cart"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockOrderRepository)
			tt.mockSetup(mockRepo)
			useCase := NewOrderUseCase(mockRepo)

			resp, err := useCase.Checkout(tt.userIDStr, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				appErr, ok := err.(*customerror.AppError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError.Type, appErr.Type)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "Jakarta", resp.ShippingAddress)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
