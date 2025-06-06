package main

import (
	"database/sql"
	"log"
	"net/http"

	"opticav2/internal/application"
	"opticav2/internal/handler"
	"opticav2/internal/infra/mysql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:rootpassword@tcp(localhost:3306)/sis_venta")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := mysql.UserRepository{DB: db}
	productRepo := mysql.ProductRepository{DB: db}

	authService := application.AuthService{Repo: userRepo}
	productService := application.ProductService{Repo: productRepo}

	authHandler := handler.AuthHandler{Service: &authService}
	productHandler := handler.ProductHandler{Service: &productService}

	mux := http.NewServeMux()
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
	mux.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
