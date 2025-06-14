package main

import (
	"database/sql"
	"log"
	"net/http"
	infraMySQL "opticav2/internal/infra/mysql"
	"opticav2/internal/application"
	"opticav2/internal/handler"

	"opticav2/internal/middleware"

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
	userRepo := infraMySQL.UserRepository{DB: gormDB}
	productRepo := infraMySQL.NewProductRepository(gormDB)
	clientRepo := infraMySQL.NewClientRepository(gormDB)
	saleRepo := infraMySQL.NewSaleRepository(gormDB, productRepo)
	paymentRepo := infraMySQL.NewPaymentRepository(gormDB)
	permissionRepo := infraMySQL.NewPermissionRepository(gormDB)
	generalLedgerRepo := infraMySQL.NewGeneralLedgerRepository(gormDB)
	configRepo := infraMySQL.NewConfigRepository(gormDB)             // Instantiate ConfigRepository
	prescriptionRepo := infraMySQL.NewPrescriptionRepository(gormDB) // Instantiate PrescriptionRepository
	permissionRepo := infraMySQL.NewPermissionRepository(gormDB) // Instantiate PermissionRepository

	// Services
	authService := application.AuthService{Repo: userRepo}
	productService := application.NewProductService(productRepo)
	userService := application.NewUserService(userRepo, permissionRepo)
	clientService := application.NewClientService(clientRepo)
	saleService := application.NewSaleService(saleRepo, paymentRepo, productRepo, clientRepo, gormDB)
	permissionService := application.NewPermissionService(permissionRepo)
	reportService := application.NewReportService(saleRepo, productRepo, generalLedgerRepo)
	statisticService := application.NewStatisticService(productRepo, saleRepo, userRepo, clientRepo)
	_ = application.NewConfigService(configRepo)
	_ = application.NewPrescriptionService(prescriptionRepo)
	// _ = application.NewPDFService(*saleService, *configService, *prescriptionService) // PDFService deshabilitado temporalmente

	// Handlers
	authHandler := handler.AuthHandler{Service: &authService}
	userHandler := handler.NewUserHandler(userService)
	clientHandler := handler.NewClientHandler(clientService)
	productHandler := handler.NewProductHandler(productService)
	saleHandler := handler.NewSaleHandler(saleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)
	reportHandler := handler.NewReportHandler(&reportService)
	statisticHandler := handler.NewStatisticHandler(&statisticService)
	// pdfHandler := handler.NewPDFHandler(pdfService) // Placeholder for PDFHandler

	// Middleware
	authzMiddleware := middleware.NewAuthorizationMiddleware(userService)

	// Router & Routes
	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", authHandler.Login)

	protectedUserHandler := authzMiddleware.RequirePermission("usuarios")(http.HandlerFunc(userHandler.HandleUserRoutes))
	mux.Handle("/api/users", protectedUserHandler)
	mux.Handle("/api/users/", protectedUserHandler)

	protectedClientHandler := authzMiddleware.RequirePermission("clientes")(http.HandlerFunc(clientHandler.HandleClientRoutes))
	mux.Handle("/api/clients", protectedClientHandler)
	mux.Handle("/api/clients/", protectedClientHandler)

	protectedProductHandler := authzMiddleware.RequirePermission("productos")(http.HandlerFunc(productHandler.HandleProductRoutes))
	mux.Handle("/api/products", protectedProductHandler)
	mux.Handle("/api/products/", protectedProductHandler)

	protectedSaleHandler := authzMiddleware.RequirePermission("ventas")(http.HandlerFunc(saleHandler.HandleSaleRoutes))
	mux.Handle("/api/sales", protectedSaleHandler)
	mux.Handle("/api/sales/", protectedSaleHandler)

	protectedPermissionHandler := authzMiddleware.RequirePermission("usuarios")(http.HandlerFunc(permissionHandler.HandlePermissionRoutes))
	mux.Handle("/api/permissions", protectedPermissionHandler)

	protectedStatsSummaryHandler := authzMiddleware.RequirePermission("estadisticas")(http.HandlerFunc(statisticHandler.GetDashboardSummary))
	mux.Handle("/api/statistics/summary", protectedStatsSummaryHandler)

	protectedLowStockHandler := authzMiddleware.RequirePermission("estadisticas")(http.HandlerFunc(statisticHandler.GetLowStockProducts))
	mux.Handle("/api/statistics/low-stock-products", protectedLowStockHandler)

	protectedTopSellingHandler := authzMiddleware.RequirePermission("estadisticas")(http.HandlerFunc(statisticHandler.GetTopSellingProducts))
	mux.Handle("/api/statistics/top-selling-products", protectedTopSellingHandler)

	protectedSalesReportHandler := authzMiddleware.RequirePermission("reporte")(http.HandlerFunc(reportHandler.GetSalesReport))
	mux.Handle("/api/reports/sales", protectedSalesReportHandler)

	// Placeholder for PDF routes
	// mux.Handle("/api/sales/{id}/receipt-pdf", authzMiddleware.RequirePermission("ventas")(http.HandlerFunc(pdfHandler.GenerateSaleReceipt)))

	log.Println("Serving static files from current working directory (expected to be project root).")
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
