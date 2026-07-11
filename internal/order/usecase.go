package order

import (
	"be-ecommerce/internal/customerror"

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
		return OrderResponse{}, customerror.NewBadRequestError("invalid user id")
	}

	// Call repository to handle transaction
	order, err := u.repo.Checkout(uid, req.ShippingAddress)
	if err != nil {
		if err.Error() == "cart is empty" {
			return OrderResponse{}, customerror.NewBadRequestError("cannot checkout an empty cart")
		}
		// In a real app we might inspect the error string to return Conflict or Bad Request if insufficient stock
		// But for now we treat other business rule violations as BadRequest if it's stock, or Internal
		return OrderResponse{}, customerror.NewInternalError("checkout transaction failed", err)
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
		return nil, customerror.NewBadRequestError("invalid user id")
	}

	orders, err := u.repo.GetMyOrders(uid)
	if err != nil {
		return nil, customerror.NewInternalError("failed to get orders", err)
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
