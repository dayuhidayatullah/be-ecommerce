package order

import (
	"be-ecommerce/internal/models"
	"github.com/google/uuid"
	"time"
)

// DTOs
type CheckoutRequest struct {
	ShippingAddress string `json:"shipping_address" validate:"required,min=10"`
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
	ImageURL  string  `json:"image_url"`
}

type OrderResponse struct {
	ID              string              `json:"id"`
	UserID          string              `json:"user_id"`
	TotalAmount     float64             `json:"total_amount"`
	Status          string              `json:"status"`
	ShippingAddress string              `json:"shipping_address"`
	CreatedAt       time.Time           `json:"created_at"`
	Items           []OrderItemResponse `json:"items"`
}

// Interfaces
type OrderRepository interface {
	Checkout(userID uuid.UUID, shippingAddress string) (*models.Order, error)
	GetMyOrders(userID uuid.UUID) ([]models.Order, error)
}

type OrderUseCase interface {
	Checkout(userID string, req CheckoutRequest) (OrderResponse, error)
	GetMyOrders(userID string) ([]OrderResponse, error)
}
