package handler

import (
	"encoding/json"
	"net/http"
	"opticav2/internal/application"
	"opticav2/internal/domain"
)

type ProductHandler struct {
	Service *application.ProductService
}

func (h ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	products, err := h.Service.List()
	if err != nil {
		http.Error(w, "error retrieving products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.Service.Create(p); err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
