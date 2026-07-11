package cart

import (
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) GetOrCreateCart(userID uuid.UUID) (*models.Cart, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Cart), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCartRepository) GetCartItem(cartID, productID uuid.UUID) (*models.CartItem, error) {
	args := m.Called(cartID, productID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.CartItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCartRepository) CreateCartItem(item *models.CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockCartRepository) UpdateCartItemQuantity(itemID uuid.UUID, addQuantity int) error {
	args := m.Called(itemID, addQuantity)
	return args.Error(0)
}

func (m *MockCartRepository) GetCartWithItems(userID uuid.UUID) (*models.Cart, []models.CartItem, error) {
	args := m.Called(userID)
	var cart *models.Cart
	if args.Get(0) != nil {
		cart = args.Get(0).(*models.Cart)
	}
	var items []models.CartItem
	if args.Get(1) != nil {
		items = args.Get(1).([]models.CartItem)
	}
	return cart, items, args.Error(2)
}
