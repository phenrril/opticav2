package main

import (
	"log"
	"net/http"

	"opticav2/internal/application"
	"opticav2/internal/handler"
	infraMySQL "opticav2/internal/infra/mysql" // Alias for clarity

	gormMySQL "gorm.io/driver/mysql" // Alias for clarity
	"gorm.io/gorm"
)

func main() {
	dsn := "root:rootpassword@tcp(localhost:3306)/sis_venta?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(gormMySQL.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Repositories
	// Assuming UserRepository might not have a New... constructor or it's a simple struct
	userRepo := infraMySQL.UserRepository{DB: gormDB}
	productRepo := infraMySQL.NewProductRepository(gormDB)
	clientRepo := infraMySQL.NewClientRepository(gormDB)
	saleRepo := infraMySQL.NewSaleRepository(gormDB, productRepo)
	paymentRepo := infraMySQL.NewPaymentRepository(gormDB)

	// Services
	// Assuming AuthService might not have a New... constructor or it's a simple struct
	authService := application.AuthService{Repo: userRepo}
	productService := application.NewProductService(productRepo)
	userService := application.NewUserService(userRepo)
	clientService := application.NewClientService(clientRepo)
	saleService := application.NewSaleService(saleRepo, paymentRepo, productRepo, clientRepo, gormDB)

	// Handlers
	// Assuming AuthHandler might not have a New... constructor or it's a simple struct
	authHandler := handler.AuthHandler{Service: &authService}
	userHandler := handler.NewUserHandler(userService)
	clientHandler := handler.NewClientHandler(clientService)
	productHandler := handler.NewProductHandler(productService)
	saleHandler := handler.NewSaleHandler(saleService)

	// Router & Routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/users", userHandler.HandleUserRoutes)
	mux.HandleFunc("/api/users/", userHandler.HandleUserRoutes)
	mux.HandleFunc("/api/clients", clientHandler.HandleClientRoutes)
	mux.HandleFunc("/api/clients/", clientHandler.HandleClientRoutes)
	mux.HandleFunc("/api/products", productHandler.HandleProductRoutes)
	mux.HandleFunc("/api/products/", productHandler.HandleProductRoutes)
	mux.HandleFunc("/api/sales", saleHandler.HandleSaleRoutes)
	mux.HandleFunc("/api/sales/", saleHandler.HandleSaleRoutes)

	// Static file serving (should be last or specific to a subpath if not root)
	log.Println("Serving static files from current working directory (expected to be project root).")
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
