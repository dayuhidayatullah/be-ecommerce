package cart

import (
	"be-ecommerce/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new instance of CartRepository
func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db}
}

func (r *cartRepository) GetOrCreateCart(userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new cart
			newCart := models.Cart{UserID: userID}
			if err := r.db.Create(&newCart).Error; err != nil {
				return nil, err
			}
			return &newCart, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetCartItem(cartID, productID uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not an error, just means it doesn't exist yet
		}
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) CreateCartItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateCartItemQuantity(itemID uuid.UUID, additionalQty int) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", itemID).
		UpdateColumn("quantity", gorm.Expr("quantity + ?", additionalQty)).Error
}

func (r *cartRepository) GetCartWithItems(userID uuid.UUID) (*models.Cart, []models.CartItem, error) {
	var cart models.Cart
	if err := r.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil // No cart yet
		}
		return nil, nil, err
	}

	var items []models.CartItem
	// Preload the Product so we can get Name and Price
	if err := r.db.Preload("Product").Where("cart_id = ?", cart.ID).Find(&items).Error; err != nil {
		return nil, nil, err
	}

	return &cart, items, nil
}
