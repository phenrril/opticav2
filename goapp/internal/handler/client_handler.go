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

type ClientHandler struct {
	ClientService *application.ClientService
}

func NewClientHandler(cs *application.ClientService) *ClientHandler {
	return &ClientHandler{ClientService: cs}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req domain.ClientCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Placeholder for authenticated user ID, as full auth middleware is not yet implemented.
	userIDPlaceholder := 1

	client, err := h.ClientService.CreateClient(req, userIDPlaceholder)
	if err != nil {
		if errors.Is(err, domain.ErrClientDNITaken) || errors.Is(err, domain.ErrClientNameTaken) {
			respondError(w, http.StatusConflict, err.Error())
		} else if strings.Contains(err.Error(), "is required") { // Basic check for validation errors
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error creating client: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusCreated, client)
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.ClientService.ListClients()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error listing clients: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request, id int) {
	client, err := h.ClientService.GetClient(id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			respondError(w, http.StatusNotFound, "Client not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Error getting client: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request, id int) {
	var req domain.ClientUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	client, err := h.ClientService.UpdateClient(id, req)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else if errors.Is(err, domain.ErrClientDNITaken) || errors.Is(err, domain.ErrClientNameTaken) {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error updating client: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) DeactivateClient(w http.ResponseWriter, r *http.Request, id int) {
	err := h.ClientService.DeactivateClient(id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error deactivating client: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Client deactivated successfully"})
	// Or respond with http.StatusNoContent and no body:
	// w.WriteHeader(http.StatusNoContent)
}

func (h *ClientHandler) ActivateClient(w http.ResponseWriter, r *http.Request, id int) {
	err := h.ClientService.ActivateClient(id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error activating client: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Client activated successfully"})
	// Or respond with http.StatusNoContent and no body:
	// w.WriteHeader(http.StatusNoContent)
}

// Master handler for /api/clients/ and /api/clients/{id}/*
func (h *ClientHandler) HandleClientRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Expected path structure: /api/clients OR /api/clients/{id} OR /api/clients/{id}/activate
	if len(pathParts) < 2 || pathParts[0] != "api" || pathParts[1] != "clients" {
		http.NotFound(w, r)
		return
	}

	if len(pathParts) == 2 { // Path is /api/clients
		switch r.Method {
		case http.MethodPost:
			h.CreateClient(w, r)
			return
		case http.MethodGet:
			h.ListClients(w, r)
			return
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/clients")
			return
		}
	} else if len(pathParts) >= 3 { // Path is /api/clients/{id} or /api/clients/{id}/action
		id, err := strconv.Atoi(pathParts[2])
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid client ID format")
			return
		}

		if len(pathParts) == 3 { // Path is /api/clients/{id}
			switch r.Method {
			case http.MethodGet:
				h.GetClient(w, r, id)
				return
			case http.MethodPut:
				h.UpdateClient(w, r, id)
				return
			case http.MethodDelete: // Using DELETE for deactivation
				h.DeactivateClient(w, r, id)
				return
			default:
				respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/clients/{id}")
				return
			}
		} else if len(pathParts) == 4 && pathParts[3] == "activate" { // Path is /api/clients/{id}/activate
			if r.Method == http.MethodPut || r.Method == http.MethodPost { // Allowing PUT or POST
				h.ActivateClient(w, r, id)
				return
			} else {
				respondError(w, http.StatusMethodNotAllowed, "Use PUT or POST for /api/clients/{id}/activate")
				return
			}
		}
	}
	http.NotFound(w, r)
}
