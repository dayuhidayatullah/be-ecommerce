package order

import (
	"be-ecommerce/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of OrderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Checkout(userID uuid.UUID, shippingAddress string) (*models.Order, error) {
	var finalOrder models.Order

	// Start Database Transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get the user's cart
		var cart models.Cart
		if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("cart is empty")
			}
			return err
		}

		// 2. Get all cart items with product details
		var cartItems []models.CartItem
		if err := tx.Preload("Product").Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
			return err
		}

		if len(cartItems) == 0 {
			return errors.New("cart is empty")
		}

		var totalAmount float64
		var orderItems []models.OrderItem

		// 3. Loop through cart items to check stock and prepare order items
		for _, item := range cartItems {
			// Check stock
			if item.Product.Stock < item.Quantity {
				return errors.New("insufficient stock for product: " + item.Product.Name)
			}

			// Deduct stock
			newStock := item.Product.Stock - item.Quantity
			if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).Update("stock", newStock).Error; err != nil {
				return err
			}

			// Calculate subtotal
			price := item.Product.Price
			totalAmount += price * float64(item.Quantity)

			// Prepare order item
			orderItems = append(orderItems, models.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     price,
			})
		}

		// 4. Create the Order record
		finalOrder = models.Order{
			UserID:          userID,
			TotalAmount:     totalAmount,
			Status:          "pending",
			ShippingAddress: shippingAddress,
		}
		if err := tx.Create(&finalOrder).Error; err != nil {
			return err
		}

		// 5. Link OrderItems to the newly created Order and save them
		for i := range orderItems {
			orderItems[i].OrderID = finalOrder.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// 6. Empty the cart (Delete all cart items)
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &finalOrder, nil
}

func (r *orderRepository) GetMyOrders(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	// Fetch orders with nested preloading for items and products
	err := r.db.Preload("User").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
