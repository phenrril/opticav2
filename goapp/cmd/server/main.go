package main

import (
	"log"
	"net/http"

	"opticav2/internal/application"
	"opticav2/internal/handler"
	infraMySQL "opticav2/internal/infra/mysql"

	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:rootpassword@tcp(localhost:3306)/sis_venta?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(gormMySQL.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	userRepo := infraMySQL.UserRepository{DB: gormDB} // Assuming direct struct instantiation or use constructor if available
	productRepo := infraMySQL.NewProductRepository(gormDB) // Use constructor
	clientRepo := infraMySQL.NewClientRepository(gormDB)   // Instantiated ClientRepository
	saleRepo := infraMySQL.NewSaleRepository(gormDB, productRepo) // Instantiate SaleRepository
	paymentRepo := infraMySQL.NewPaymentRepository(gormDB)     // Instantiate PaymentRepository

	authService := application.AuthService{Repo: userRepo}    // Assuming direct struct instantiation or use constructor
	productService := application.NewProductService(productRepo) // Use constructor
	userService := application.NewUserService(userRepo)       // Instantiated UserService
	clientService := application.NewClientService(clientRepo)   // Instantiated ClientService
	saleService := application.NewSaleService(saleRepo, paymentRepo, productRepo, clientRepo, gormDB) // Instantiate SaleService

	authHandler := handler.AuthHandler{Service: &authService} // Assuming direct struct instantiation or use constructor
	userHandler := handler.NewUserHandler(userService)           // UserHandler is now created
	clientHandler := handler.NewClientHandler(clientService)     // ClientHandler is now created
	productHandler := handler.NewProductHandler(productService)  // Use constructor for ProductHandler
	saleHandler := handler.NewSaleHandler(saleService)        // SaleHandler is now created

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", userHandler.HandleUserRoutes)    // Handles POST for create, GET for list
	mux.HandleFunc("/api/users/", userHandler.HandleUserRoutes)   // Handles GET /id, PUT /id, DELETE /id, PUT /id/activate
	mux.HandleFunc("/api/clients", clientHandler.HandleClientRoutes)  // Handles POST for create, GET for list
	mux.HandleFunc("/api/clients/", clientHandler.HandleClientRoutes) // Handles GET /id, PUT /id, DELETE /id, PUT /id/activate
	mux.HandleFunc("/api/products", productHandler.HandleProductRoutes)  // Handles all product routes
	mux.HandleFunc("/api/products/", productHandler.HandleProductRoutes) // Handles all product routes including /id/*
	mux.HandleFunc("/api/sales", saleHandler.HandleSaleRoutes)    // Handles POST for create, GET for list
	mux.HandleFunc("/api/sales/", saleHandler.HandleSaleRoutes)   // Handles GET /id, POST /id/payments, DELETE /id
	mux.HandleFunc("/api/login", authHandler.Login)
	log.Println("Serving static files from current working directory (expected to be project root).")
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
