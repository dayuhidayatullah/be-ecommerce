package database

import (
	"log"

	"be-ecommerce/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database successfully.")

	// Enable uuid-ossp extension for uuid_generate_v4()
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Auto Migrate the models
	log.Println("Running database migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return nil, err
	}
	log.Println("Database migrations completed.")

	return db, nil
}
