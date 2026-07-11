package cart

import (
	"be-ecommerce/internal/models"
	"github.com/google/uuid"
)

// DTOs
type AddToCartRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type CartItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type CartResponse struct {
	ID          string             `json:"id"`
	UserID      string             `json:"user_id"`
	Items       []CartItemResponse `json:"items"`
	TotalAmount float64            `json:"total_amount"`
}

// Interfaces
type CartRepository interface {
	GetOrCreateCart(userID uuid.UUID) (*models.Cart, error)
	GetCartItem(cartID, productID uuid.UUID) (*models.CartItem, error)
	CreateCartItem(item *models.CartItem) error
	UpdateCartItemQuantity(itemID uuid.UUID, additionalQty int) error
	GetCartWithItems(userID uuid.UUID) (*models.Cart, []models.CartItem, error)
}

type CartUseCase interface {
	AddToCart(userID string, req AddToCartRequest) error
	GetMyCart(userID string) (CartResponse, error)
}
