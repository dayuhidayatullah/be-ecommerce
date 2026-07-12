package cart

import (
	"be-ecommerce/internal/customerror"
	"be-ecommerce/internal/models"

	"github.com/google/uuid"
)

type cartUseCase struct {
	repo CartRepository
}

// NewCartUseCase creates a new instance of CartUseCase
func NewCartUseCase(repo CartRepository) CartUseCase {
	return &cartUseCase{repo: repo}
}

func (u *cartUseCase) AddToCart(userID string, req AddToCartRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return customerror.NewBadRequestError("invalid user id")
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return customerror.NewBadRequestError("invalid product id")
	}

	// 1. Get or create cart for user
	cart, err := u.repo.GetOrCreateCart(uid)
	if err != nil {
		return customerror.NewInternalError("failed to get or create cart", err)
	}

	// 2. Check if product already exists in cart
	existingItem, err := u.repo.GetCartItem(cart.ID, productID)
	if err != nil {
		return customerror.NewInternalError("failed to get cart item", err)
	}

	// 3. If exists, update quantity
	if existingItem != nil {
		if err := u.repo.UpdateCartItemQuantity(existingItem.ID, req.Quantity); err != nil {
			return customerror.NewInternalError("failed to update cart item quantity", err)
		}
		return nil
	}

	// 4. If not exists, create new cart item
	newItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  req.Quantity,
	}
	
	err = u.repo.CreateCartItem(newItem)
	if err != nil {
		return customerror.NewNotFoundError("failed to add item, ensure product exists")
	}

	return nil
}

func (u *cartUseCase) GetMyCart(userID string) (CartResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return CartResponse{}, customerror.NewBadRequestError("invalid user id")
	}

	cart, items, err := u.repo.GetCartWithItems(uid)
	if err != nil {
		return CartResponse{}, customerror.NewInternalError("failed to get cart with items", err)
	}

	// Handle case where user hasn't added anything yet
	if cart == nil {
		return CartResponse{
			UserID: userID,
			Items:  []CartItemResponse{},
		}, nil
	}

	var itemResponses []CartItemResponse
	var totalAmount float64

	for _, item := range items {
		subtotal := item.Product.Price * float64(item.Quantity)
		totalAmount += subtotal

		itemResponses = append(itemResponses, CartItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Name:      item.Product.Name,
			Price:     item.Product.Price,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
			ImageURL:  item.Product.ImageURL,
		})
	}

	// Avoid nil slice in JSON
	if itemResponses == nil {
		itemResponses = []CartItemResponse{}
	}

	return CartResponse{
		ID:          cart.ID.String(),
		UserID:      cart.UserID.String(),
		Items:       itemResponses,
		TotalAmount: totalAmount,
	}, nil
}
