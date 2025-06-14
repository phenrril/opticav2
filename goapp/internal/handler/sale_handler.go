package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv" // For parsing ID from URL
	"strings" // For routing logic

	"opticav2/internal/application"
	"opticav2/internal/domain"
)

type SaleHandler struct {
	SaleService *application.SaleService
}

func NewSaleHandler(ss *application.SaleService) *SaleHandler {
	return &SaleHandler{SaleService: ss}
}

func (h *SaleHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	userIDPlaceholder := int(1) // Simulate authenticated user ID

	sale, err := h.SaleService.CreateSale(req, userIDPlaceholder)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSaleData) || strings.Contains(err.Error(), "not found for sale item") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, domain.ErrInsufficientStock) {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error creating sale: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusCreated, sale)
}

func (h *SaleHandler) ListSales(w http.ResponseWriter, r *http.Request) {
	userIDPlaceholder := int(1) // Simulate authenticated user ID for filtering or context

	// Basic filter parsing from query parameters (example)
	filters := make(map[string]interface{})
	if clientIDStr := r.URL.Query().Get("client_id"); clientIDStr != "" {
		if clientID, err := strconv.ParseUint(clientIDStr, 10, 32); err == nil {
			filters["client_id"] = uint(clientID)
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}
	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom // Assuming YYYY-MM-DD format
	}
	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		filters["date_to"] = dateTo // Assuming YYYY-MM-DD format
	}
	// For role-based filtering, service might apply user_id filter by default
	// filters["user_id"] = userIDPlaceholder // If service doesn't automatically filter by user

	sales, err := h.SaleService.ListSales(userIDPlaceholder, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error listing sales: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, sales)
}

func (h *SaleHandler) GetSale(w http.ResponseWriter, r *http.Request, id int) {
	userIDPlaceholder := uint(1) // Simulate authenticated user ID for context/auth check in service

	sale, err := h.SaleService.GetSale(id, userIDPlaceholder)
	if err != nil {
		if errors.Is(err, domain.ErrSaleNotFound) {
			respondError(w, http.StatusNotFound, "Sale not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Error getting sale: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, sale)
}

func (h *SaleHandler) AddPaymentToSale(w http.ResponseWriter, r *http.Request, id int) {
	var req domain.AddPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	processedByUserIDPlaceholder := int(1) // Simulate authenticated user ID

	req.SaleID = id // Ensure SaleID from path is used

	sale, err := h.SaleService.AddPaymentToSale(id, req, processedByUserIDPlaceholder)
	if err != nil {
		if errors.Is(err, domain.ErrSaleNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "must be positive") || strings.Contains(err.Error(), "cannot add payment to sale with status") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error adding payment: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, sale)
}

func (h *SaleHandler) CancelSale(w http.ResponseWriter, r *http.Request, id int) {
	cancelledByUserIDPlaceholder := int(1) // Simulate authenticated user ID

	sale, err := h.SaleService.CancelSale(id, cancelledByUserIDPlaceholder)
	if err != nil {
		if errors.Is(err, domain.ErrSaleNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "already cancelled") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error cancelling sale: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, sale) // Return updated sale with "Cancelled" status
}

// HandleSaleRoutes is the master handler for /api/sales/*
func (h *SaleHandler) HandleSaleRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Expected: /api/sales OR /api/sales/{id} OR /api/sales/{id}/action
	if len(pathParts) < 2 || pathParts[0] != "api" || pathParts[1] != "sales" {
		http.NotFound(w, r)
		return
	}

	if len(pathParts) == 2 { // Path is /api/sales
		switch r.Method {
		case http.MethodPost:
			h.CreateSale(w, r)
			return
		case http.MethodGet:
			h.ListSales(w, r)
			return
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/sales")
			return
		}
	} else if len(pathParts) >= 3 { // Path is /api/sales/{id} or /api/sales/{id}/action
		id64, err := strconv.ParseUint(pathParts[2], 10, 32)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid sale ID format")
			return
		}
		id := int(id64)

		if len(pathParts) == 3 { // Path is /api/sales/{id}
			switch r.Method {
			case http.MethodGet:
				h.GetSale(w, r, id)
				return
			case http.MethodDelete: // For cancelling a sale
				h.CancelSale(w, r, id)
				return
			default:
				respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/sales/{id}")
				return
			}
		} else if len(pathParts) == 4 { // Path is /api/sales/{id}/action
			action := pathParts[3]
			switch action {
			case "payments":
				if r.Method == http.MethodPost {
					h.AddPaymentToSale(w, r, id)
					return
				} else {
					respondError(w, http.StatusMethodNotAllowed, "Use POST for /api/sales/{id}/payments")
					return
				}
			// Add other actions like "refund", "resend_receipt" if needed later
			default:
				http.NotFound(w, r)
				return
			}
		}
	}
	http.NotFound(w, r)
}
