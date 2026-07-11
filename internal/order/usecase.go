package order

import (
	"errors"

	"github.com/google/uuid"
)

type orderUseCase struct {
	repo OrderRepository
}

// NewOrderUseCase creates a new instance of OrderUseCase
func NewOrderUseCase(repo OrderRepository) OrderUseCase {
	return &orderUseCase{repo: repo}
}

func (u *orderUseCase) Checkout(userID string, req CheckoutRequest) (OrderResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return OrderResponse{}, errors.New("invalid user id")
	}

	// Call repository to handle transaction
	order, err := u.repo.Checkout(uid, req.ShippingAddress)
	if err != nil {
		return OrderResponse{}, err
	}

	return OrderResponse{
		ID:              order.ID.String(),
		UserID:          order.UserID.String(),
		TotalAmount:     order.TotalAmount,
		Status:          order.Status,
		ShippingAddress: order.ShippingAddress,
		CreatedAt:       order.CreatedAt,
		Items:           []OrderItemResponse{},
	}, nil
}

func (u *orderUseCase) GetMyOrders(userID string) ([]OrderResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	orders, err := u.repo.GetMyOrders(uid)
	if err != nil {
		return nil, err
	}

	var responses []OrderResponse
	for _, o := range orders {
		var itemResps []OrderItemResponse
		for _, item := range o.OrderItems {
			itemResps = append(itemResps, OrderItemResponse{
				ID:        item.ID.String(),
				ProductID: item.ProductID.String(),
				Name:      item.Product.Name,
				Quantity:  item.Quantity,
				Price:     item.Price,
				Subtotal:  item.Price * float64(item.Quantity),
			})
		}
		
		responses = append(responses, OrderResponse{
			ID:              o.ID.String(),
			UserID:          o.UserID.String(),
			TotalAmount:     o.TotalAmount,
			Status:          o.Status,
			ShippingAddress: o.ShippingAddress,
			CreatedAt:       o.CreatedAt,
			Items:           itemResps, 
		})
	}
	return responses, nil
}
