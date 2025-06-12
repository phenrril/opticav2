package handler

import (
	"encoding/json"
	"net/http"
	"opticav2/internal/application"
	// "opticav2/internal/domain" // Not strictly needed for this handler if service returns concrete types
	"time"
	// "strconv" // Not needed for this specific handler's current requirements
	// "strings" // Not needed for this specific handler's current requirements
)

type ReportHandler struct {
	ReportService *application.ReportService
}

func NewReportHandler(rs *application.ReportService) *ReportHandler {
	return &ReportHandler{ReportService: rs}
}

// Helper to respond with JSON - Copied
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

// Helper to respond with an error - Copied
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// GetSalesReport handles GET /api/reports/sales
func (h *ReportHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from_date") // Expected format YYYY-MM-DD
	toStr := r.URL.Query().Get("to_date")     // Expected format YYYY-MM-DD

	if fromStr == "" || toStr == "" {
		respondError(w, http.StatusBadRequest, "from_date and to_date query parameters are required (YYYY-MM-DD)")
		return
	}

	layout := "2006-01-02"
	fromDate, err1 := time.Parse(layout, fromStr)
	toDate, err2 := time.Parse(layout, toStr)

	if err1 != nil || err2 != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format. Please use YYYY-MM-DD.")
		return
	}

	// Ensure toDate covers the entire day if the backend query uses BETWEEN or <= toDate.
	// For instance, if toDate is "2023-10-20", it means up to "2023-10-20 00:00:00".
	// To include the whole day, it should be "2023-10-20 23:59:59".
	// The ReportService already handles this by adding 23h59m59s to toDate.

	// Placeholder for userID from context (if needed for filtering by user in service)
	// For this report, the PHP version seems system-wide beyond permission check.
	// So, userID for filtering might be 0 or a special value if service expects it.
	// The service method `GenerateSalesReport` has a userID param.
	userIDPlaceholder := uint(0) // 0 could mean "all users" or "system context" for the report
	// If a specific user context is absolutely required from middleware, fetch it here.
	// For now, assuming the report itself is not filtered by the requesting user unless specified in `otherFilters`.

	// Placeholder for other potential filters from query params
	otherFilters := make(map[string]interface{})
	// Example: if r.URL.Query().Get("some_filter") != "" {
	//     otherFilters["some_key"] = r.URL.Query().Get("some_filter")
	// }

	report, err := (*h.ReportService).GenerateSalesReport(fromDate, toDate, userIDPlaceholder, otherFilters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error generating sales report: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, report)
}

// HandleReportRoutes could be used if more report endpoints are added under /api/reports/*
/*
func (h *ReportHandler) HandleReportRoutes(w http.ResponseWriter, r *http.Request) {
    // Example routing if needed
    if (r.URL.Path == "/api/reports/sales" || r.URL.Path == "/api/reports/sales/") && r.Method == http.MethodGet {
        h.GetSalesReport(w, r)
        return
    }
    http.NotFound(w, r)
}
*/
