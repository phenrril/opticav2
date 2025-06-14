package handler

import (
	"opticav2/internal/application"
	// "strconv" // For parsing query params later if needed
	// "opticav2/internal/domain" // Not strictly needed if service returns concrete types
)

type StatisticHandler struct {
	StatisticService *application.StatisticService
}

func NewStatisticHandler(ss *application.StatisticService) *StatisticHandler {
	return &StatisticHandler{StatisticService: ss}
}

/*
// GetDashboardSummary handles GET /api/statistics/summary
func (h *StatisticHandler) GetDashboardSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.StatisticService.GetDashboardSummaryCounts()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching dashboard summary: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, summary)
}


// GetLowStockProducts handles GET /api/statistics/low-stock-products
func (h *StatisticHandler) GetLowStockProducts(w http.ResponseWriter, r *http.Request) {
	// Fixed threshold for now, as per PHP chart.php (value often 5 or 10)
	// The PHP query used `existencia <= 5`. Let's use 10 as per task example.
	threshold := 10

	// Optional: Parse query param for threshold
	// if thresholdStr := r.URL.Query().Get("threshold"); thresholdStr != "" {
	// 	if val, err := strconv.Atoi(thresholdStr); err == nil && val > 0 {
	// 		threshold = val
	// 	}
	// }

	products, err := h.StatisticService.GetLowStockProducts(threshold)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching low stock products: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, products)
}

// GetTopSellingProducts handles GET /api/statistics/top-selling-products
func (h *StatisticHandler) GetTopSellingProducts(w http.ResponseWriter, r *http.Request) {
	// Fixed limit for now, as per PHP chart.php (value often 5)
	limit := 5

	// Optional: Parse query params for limit, from_date, to_date
	// if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
	// 	if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
	// 		limit = val
	// 	}
	// }
	// fromDateStr := r.URL.Query().Get("from_date")
	// toDateStr := r.URL.Query().Get("to_date")

	// Using empty strings for dates means no date filtering in the service for now.
	products, err := h.StatisticService.GetTopSellingProducts(limit, "", "")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching top selling products: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, products)
}

// HandleStatisticRoutes could be used if a master router for /api/statistics/* is preferred later.
// For now, individual handlers are registered in main.go.
/*
func (h *StatisticHandler) HandleStatisticRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/statistics")
	switch path {
	case "/summary":
		if r.Method == http.MethodGet {
			h.GetDashboardSummary(w, r)
			return
		}
	case "/low-stock-products":
		if r.Method == http.MethodGet {
			h.GetLowStockProducts(w, r)
			return
		}
	case "/top-selling-products":
		if r.Method == http.MethodGet {
			h.GetTopSellingProducts(w, r)
			return
		}
	default:
		http.NotFound(w,r)
		return
	}
	respondError(w, http.StatusMethodNotAllowed, "Method not allowed for this statistics endpoint")
}
*/
