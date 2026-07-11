package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"type:varchar(20);default:'customer'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"uniqueIndex;not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null"`
	Category    Category  `gorm:"foreignKey:CategoryID"`
	Name        string    `gorm:"not null"`
	Description string
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Stock       int       `gorm:"not null;default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Cart      Cart      `gorm:"foreignKey:CartID"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	Quantity  int       `gorm:"not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	User            User      `gorm:"foreignKey:UserID"`
	TotalAmount     float64   `gorm:"type:decimal(10,2);not null"`
	Status          string    `gorm:"type:varchar(50);default:'pending'"`
	ShippingAddress string    `gorm:"type:text;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Order     Order     `gorm:"foreignKey:OrderID"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"type:decimal(10,2);not null"` // Price at checkout
	CreatedAt time.Time
	UpdatedAt time.Time
}
