package main

import (
	"log"
	"net/http"

	"opticav2/internal/application"
	"opticav2/internal/handler"

	"opticav2/internal/middleware"            // Import middleware package
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

	permissionRepo := infraMySQL.NewPermissionRepository(gormDB) // Instantiate PermissionRepository


	// Services
	// Assuming AuthService might not have a New... constructor or it's a simple struct
	authService := application.AuthService{Repo: userRepo}
	productService := application.NewProductService(productRepo)

	userService := application.NewUserService(userRepo, permissionRepo) // Update UserService instantiation
	clientService := application.NewClientService(clientRepo)
	saleService := application.NewSaleService(saleRepo, paymentRepo, productRepo, clientRepo, gormDB)
	permissionService := application.NewPermissionService(permissionRepo) // Instantiate PermissionService


	// Handlers
	// Assuming AuthHandler might not have a New... constructor or it's a simple struct
	authHandler := handler.AuthHandler{Service: &authService}
	userHandler := handler.NewUserHandler(userService)
	clientHandler := handler.NewClientHandler(clientService)
	productHandler := handler.NewProductHandler(productService)
	saleHandler := handler.NewSaleHandler(saleService)

	permissionHandler := handler.NewPermissionHandler(permissionService)

	// Middleware
	authzMiddleware := middleware.NewAuthorizationMiddleware(userService)


	// Router & Routes
	mux := http.NewServeMux()

	// API routes

	mux.HandleFunc("/api/login", authHandler.Login) // Login usually doesn't need authz

	// Protected User Routes
	protectedUserHandler := authzMiddleware.RequirePermission("usuarios")(http.HandlerFunc(userHandler.HandleUserRoutes))
	mux.Handle("/api/users", protectedUserHandler)
	mux.Handle("/api/users/", protectedUserHandler)

	// Protected Client Routes
	protectedClientHandler := authzMiddleware.RequirePermission("clientes")(http.HandlerFunc(clientHandler.HandleClientRoutes))
	mux.Handle("/api/clients", protectedClientHandler)
	mux.Handle("/api/clients/", protectedClientHandler)

	// Protected Product Routes
	protectedProductHandler := authzMiddleware.RequirePermission("productos")(http.HandlerFunc(productHandler.HandleProductRoutes))
	mux.Handle("/api/products", protectedProductHandler)
	mux.Handle("/api/products/", protectedProductHandler)

	// Protected Sales Routes
	// Assuming "ventas" for general sales access. Granular permissions (e.g. "nueva_venta") would require more complex routing or checks within handler.
	protectedSaleHandler := authzMiddleware.RequirePermission("ventas")(http.HandlerFunc(saleHandler.HandleSaleRoutes))
	mux.Handle("/api/sales", protectedSaleHandler)
	mux.Handle("/api/sales/", protectedSaleHandler)

	// Permissions listing route (could also be protected, e.g., by "usuarios" or a specific "view_permissions" perm)
	// For now, let's assume it's protected by "usuarios" permission as it's related to user roles.
	protectedPermissionHandler := authzMiddleware.RequirePermission("usuarios")(http.HandlerFunc(permissionHandler.HandlePermissionRoutes))
	mux.Handle("/api/permissions", protectedPermissionHandler)


	// Static file serving (should be last or specific to a subpath if not root)
	log.Println("Serving static files from current working directory (expected to be project root).")
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
