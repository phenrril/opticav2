package handler

import (
	"net/http"
	"opticav2/internal/application"
	// "opticav2/internal/domain" // Not strictly needed for this handler if only listing all
	// "strconv" // Not needed for just /api/permissions
	// "strings" // Not needed for just /api/permissions
)

type PermissionHandler struct {
	PermissionService *application.PermissionService
}

func NewPermissionHandler(ps *application.PermissionService) *PermissionHandler {
	return &PermissionHandler{PermissionService: ps}
}

func (h *PermissionHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.PermissionService.ListAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error listing permissions: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, permissions)
}

// HandlePermissionRoutes is the master handler for /api/permissions/*
func (h *PermissionHandler) HandlePermissionRoutes(w http.ResponseWriter, r *http.Request) {
	// path := strings.TrimSpace(r.URL.Path) // For more complex routing under /api/permissions if needed later

	// Only /api/permissions (GET) is defined for now.
	if r.URL.Path == "/api/permissions" || r.URL.Path == "/api/permissions/" {
		if r.Method == http.MethodGet {
			h.ListPermissions(w, r)
			return
		} else {
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/permissions")
			return
		}
	}
	http.NotFound(w, r)
}
