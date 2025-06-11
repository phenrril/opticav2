package handler

import (
	"encoding/json"
	"net/http"
	"opticav2/internal/application"
	"opticav2/internal/domain"
	"strconv" // For parsing ID from URL
	"strings" // For routing logic
)

type UserHandler struct {
	UserService *application.UserService
}

func NewUserHandler(us *application.UserService) *UserHandler {
	return &UserHandler{UserService: us}
}

// Helper to respond with JSON
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

// Helper to respond with an error
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req domain.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := h.UserService.CreateUser(req)
	if err != nil {
		// Basic error handling, can be more specific
		if strings.Contains(err.Error(), "already exists") {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error creating user: "+err.Error())
		}
		return
	}
	// Exclude password from response (already handled by json:"-" in User struct)
	respondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.ListUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error listing users: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, id int) {
	user, err := h.UserService.GetUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") { // Or check for gorm.ErrRecordNotFound type from service
			respondError(w, http.StatusNotFound, "User not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Error getting user: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var req domain.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := h.UserService.UpdateUser(id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "already in use") {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error updating user: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request, id int) {
	err := h.UserService.DeactivateUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error deactivating user: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "User deactivated successfully"})
}

func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request, id int) {
	err := h.UserService.ActivateUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error activating user: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "User activated successfully"})
}

// Master handler for /api/users/ and /api/users/{id}/*
func (h *UserHandler) HandleUserRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// Expected path structure: /api/users OR /api/users/{id} OR /api/users/{id}/activate
	if len(pathParts) < 2 || pathParts[0] != "api" || pathParts[1] != "users" {
		http.NotFound(w, r)
		return
	}

	if len(pathParts) == 2 { // Path is /api/users
		switch r.Method {
		case http.MethodPost:
			h.CreateUser(w, r)
			return
		case http.MethodGet:
			h.ListUsers(w, r)
			return
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/users")
			return
		}
	} else if len(pathParts) >= 3 { // Path is /api/users/{id} or /api/users/{id}/action
		id, err := strconv.Atoi(pathParts[2])
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID format")
			return
		}

		if len(pathParts) == 3 { // Path is /api/users/{id}
			switch r.Method {
			case http.MethodGet:
				h.GetUser(w, r, id)
				return
			case http.MethodPut:
				h.UpdateUser(w, r, id)
				return
			case http.MethodDelete: // Using DELETE for deactivation as per RESTful conventions
				h.DeactivateUser(w, r, id)
				return
			default:
				respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/users/{id}")
				return
			}
		} else if len(pathParts) == 4 && pathParts[3] == "activate" { // Path is /api/users/{id}/activate
			if r.Method == http.MethodPut || r.Method == http.MethodPost { // Allowing PUT or POST for activation
				h.ActivateUser(w, r, id)
				return
			} else {
				respondError(w, http.StatusMethodNotAllowed, "Use PUT or POST for /api/users/{id}/activate")
				return
			}
		}
	}
	http.NotFound(w, r)
}
