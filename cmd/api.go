package main

import (
	"net/http"
	"time"

	"be-ecommerce/internal/auth"
	"be-ecommerce/internal/middleware"
	"be-ecommerce/internal/product"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type application struct {
	config   config
	db       *gorm.DB
	validate *validator.Validate
}
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

  // A good base middleware stack
  r.Use(chimiddleware.RequestID)
  r.Use(chimiddleware.RealIP)
  r.Use(chimiddleware.Logger)
  r.Use(chimiddleware.Recoverer)

  // Set a timeout value on the request context (ctx)
  r.Use(chimiddleware.Timeout(60 * time.Second))

  r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hi"))
  })

  r.Route("/api/v1", func(r chi.Router) {
      // --- Dependency Injection ---
      
      // Auth module
      authRepo := auth.NewAuthRepository(app.db)
      authUseCase := auth.NewAuthUseCase(authRepo)
      authHandler := auth.NewHttpHandler(authUseCase, app.validate)
      
      // Product module
      productRepo := product.NewProductRepository(app.db)
      productUseCase := product.NewProductUseCase(productRepo)
      productHandler := product.NewHttpHandler(productUseCase, app.validate)

      // --- Public Routes ---
      r.Post("/register", authHandler.Register)
      r.Post("/login", authHandler.Login)
      
      r.Get("/categories", productHandler.GetCategories)
      r.Get("/products", productHandler.GetProducts)

      // --- Protected Routes (Require Login) ---
      r.Group(func(r chi.Router) {
          r.Use(middleware.RequireAuth)
          
          // --- Admin Only Routes ---
          r.Group(func(r chi.Router) {
              r.Use(middleware.RequireRole("admin"))
              
              r.Post("/categories", productHandler.CreateCategory)
              r.Post("/products", productHandler.CreateProduct)
          })
      })
  })

  return r
}
