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

	userRepo := infraMySQL.UserRepository{DB: gormDB}
	productRepo := infraMySQL.ProductRepository{DB: gormDB}
	clientRepo := infraMySQL.NewClientRepository(gormDB) // Instantiated ClientRepository

	authService := application.AuthService{Repo: userRepo}
	productService := application.ProductService{Repo: productRepo}
	userService := application.NewUserService(userRepo)     // Instantiated UserService
	clientService := application.NewClientService(clientRepo) // Instantiated ClientService

	authHandler := handler.AuthHandler{Service: &authService}
	productHandler := handler.ProductHandler{Service: &productService}
	userHandler := handler.NewUserHandler(userService)         // UserHandler is now created
	clientHandler := handler.NewClientHandler(clientService)   // ClientHandler is now created

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users", userHandler.HandleUserRoutes)    // Handles POST for create, GET for list
	mux.HandleFunc("/api/users/", userHandler.HandleUserRoutes)   // Handles GET /id, PUT /id, DELETE /id, PUT /id/activate
	mux.HandleFunc("/api/clients", clientHandler.HandleClientRoutes)  // Handles POST for create, GET for list
	mux.HandleFunc("/api/clients/", clientHandler.HandleClientRoutes) // Handles GET /id, PUT /id, DELETE /id, PUT /id/activate
	mux.HandleFunc("/api/login", authHandler.Login)
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.List(w, r)
		case http.MethodPost:
			productHandler.Create(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	log.Println("Serving static files from current working directory (expected to be project root).")
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
