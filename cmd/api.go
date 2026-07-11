package main

import (
	"net/http"
	"time"

	"be-ecommerce/internal/auth"
	"be-ecommerce/internal/cart"
	"be-ecommerce/internal/middleware"
	"be-ecommerce/internal/order"
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

      // Cart module
      cartRepo := cart.NewCartRepository(app.db)
      cartUseCase := cart.NewCartUseCase(cartRepo)
      cartHandler := cart.NewHttpHandler(cartUseCase, app.validate)
      
      // Order module
      orderRepo := order.NewOrderRepository(app.db)
      orderUseCase := order.NewOrderUseCase(orderRepo)
      orderHandler := order.NewHttpHandler(orderUseCase, app.validate)

      // --- Public Routes ---
      r.Post("/register", authHandler.Register)
      r.Post("/login", authHandler.Login)
      
      r.Get("/categories", productHandler.GetCategories)
      r.Get("/products", productHandler.GetProducts)

      // --- Protected Routes (Require Login) ---
      r.Group(func(r chi.Router) {
          r.Use(middleware.RequireAuth)
          
          // Cart routes (User)
          r.Get("/cart", cartHandler.GetMyCart)
          r.Post("/cart/items", cartHandler.AddToCart)
          
          // Order routes (User)
          r.Get("/orders", orderHandler.GetMyOrders)
          r.Post("/orders/checkout", orderHandler.Checkout)

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
