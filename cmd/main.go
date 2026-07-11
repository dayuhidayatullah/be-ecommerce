package main

import (
	"log"
	"net/http"
	"os"

	"be-ecommerce/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	cfg := config{
		addr: ":" + port,
		db: dbConfig{
			dsn: dbDsn,
		},
	}

	// Initialize database
	db, err := database.New(cfg.db.dsn)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	// Initialize validator
	validate := validator.New()

	api := &application{
		config:   cfg,
		db:       db,
		validate: validate,
	}

	if err := api.run(api.mount()); err != nil {
		log.Printf("Error starting server: %v", err)
		os.Exit(1)
	}

}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}

	log.Printf("Server is starting on %s", app.config.addr)
	return srv.ListenAndServe()
}
