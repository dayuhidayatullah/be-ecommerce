package order

import (
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Checkout(userID uuid.UUID, shippingAddress string) (*models.Order, error) {
	args := m.Called(userID, shippingAddress)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderRepository) GetMyOrders(userID uuid.UUID) ([]models.Order, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}
