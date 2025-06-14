package handler

import (
	"encoding/json"
	"errors" // For domain error checking
	"net/http"
	"strconv" // For parsing ID from URL
	"strings" // For routing logic

	"opticav2/internal/application"
	"opticav2/internal/domain"
)

type ProductHandler struct {
	ProductService *application.ProductService
}

func NewProductHandler(ps *application.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: ps}
}

// Helper to respond with JSON - Copied for now
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Helper to respond with an error - Copied for now
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// Create handles POST requests to /api/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Placeholder for authenticated user ID
	userIDPlaceholder := 1

	product, err := h.ProductService.CreateProduct(req, userIDPlaceholder)
	if err != nil {
		if errors.Is(err, domain.ErrProductCodeTaken) {
			respondError(w, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "is required") { // Basic check for validation errors
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error creating product: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

// List handles GET requests to /api/products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.ListProducts()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error listing products: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request, id int) {
	product, err := h.ProductService.GetProduct(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			respondError(w, http.StatusNotFound, "Product not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Error getting product: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, id int) {
	var req domain.ProductUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	product, err := h.ProductService.UpdateProduct(id, req)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else if errors.Is(err, domain.ErrProductCodeTaken) {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error updating product: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProductStock(w http.ResponseWriter, r *http.Request, id int) {
	var req domain.ProductStockUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	product, err := h.ProductService.UpdateProductStock(id, req)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "stock cannot be negative") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error updating product stock: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) DeactivateProduct(w http.ResponseWriter, r *http.Request, id int) {
	err := h.ProductService.DeactivateProduct(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error deactivating product: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Product deactivated successfully"})
}

func (h *ProductHandler) ActivateProduct(w http.ResponseWriter, r *http.Request, id int) {
	err := h.ProductService.ActivateProduct(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error activating product: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Product activated successfully"})
}

// HandleProductRoutes is the master handler for /api/products/*
func (h *ProductHandler) HandleProductRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Expected: /api/products OR /api/products/{id} OR /api/products/{id}/action
	if len(pathParts) < 2 || pathParts[0] != "api" || pathParts[1] != "products" {
		http.NotFound(w, r)
		return
	}

	if len(pathParts) == 2 { // Path is /api/products
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r) // Renamed from CreateProduct to Create for consistency
			return
		case http.MethodGet:
			h.List(w, r) // Renamed from ListProducts to List
			return
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/products")
			return
		}
	} else if len(pathParts) >= 3 { // Path is /api/products/{id} or /api/products/{id}/action
		id, err := strconv.Atoi(pathParts[2])
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid product ID format")
			return
		}

		if len(pathParts) == 3 { // Path is /api/products/{id}
			switch r.Method {
			case http.MethodGet:
				h.GetProduct(w, r, id)
				return
			case http.MethodPut: // General update
				h.UpdateProduct(w, r, id)
				return
			case http.MethodDelete: // Deactivation
				h.DeactivateProduct(w, r, id)
				return
			default:
				respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/products/{id}")
				return
			}
		} else if len(pathParts) == 4 { // Path is /api/products/{id}/action
			action := pathParts[3]
			switch action {
			case "stock":
				if r.Method == http.MethodPut || r.Method == http.MethodPost {
					h.UpdateProductStock(w, r, id)
					return
				} else {
					respondError(w, http.StatusMethodNotAllowed, "Use PUT or POST for /api/products/{id}/stock")
					return
				}
			case "activate":
				if r.Method == http.MethodPut || r.Method == http.MethodPost {
					h.ActivateProduct(w, r, id)
					return
				} else {
					respondError(w, http.StatusMethodNotAllowed, "Use PUT or POST for /api/products/{id}/activate")
					return
				}
			default:
				http.NotFound(w, r)
				return
			}
		}
	}
	http.NotFound(w, r)
}
