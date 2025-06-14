package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"opticav2/internal/application"
	"opticav2/internal/domain"
)

type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

type UserHandler struct {
	UserService *application.UserService
}

func NewUserHandler(us *application.UserService) *UserHandler {
	return &UserHandler{UserService: us}
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
		if strings.Contains(err.Error(), "already exists") {
			respondError(w, http.StatusConflict, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error creating user: "+err.Error())
		}
		return
	}
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, id uint) {
	user, err := h.UserService.GetUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "User not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Error getting user: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id uint) {
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

func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request, id uint) {
	err := h.UserService.DeactivateUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "User not found for deactivation")
		} else {
			respondError(w, http.StatusInternalServerError, "Error deactivating user: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "User deactivated successfully"})
}

func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request, id uint) {
	err := h.UserService.ActivateUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "User not found for activation")
		} else {
			respondError(w, http.StatusInternalServerError, "Error activating user: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "User activated successfully"})
}

func (h *UserHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request, userID uint) {
	permissions, err := h.UserService.GetUserPermissions(userID)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error getting user permissions: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, permissions)
}

func (h *UserHandler) AssignPermissionsToUser(w http.ResponseWriter, r *http.Request, userID uint) {
	var req AssignPermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	err := h.UserService.AssignPermissionsToUser(userID, req.PermissionIDs)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			respondError(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "permission IDs are invalid") {
			respondError(w, http.StatusBadRequest, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Error assigning permissions: "+err.Error())
		}
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Permissions assigned successfully"})
}

func (h *UserHandler) HandleUserRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(pathParts) < 2 || pathParts[0] != "api" || pathParts[1] != "users" {
		http.NotFound(w, r)
		return
	}

	if len(pathParts) == 2 {
		switch r.Method {
		case http.MethodPost:
			h.CreateUser(w, r)
		case http.MethodGet:
			h.ListUsers(w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/users")
		}
		return
	} else if len(pathParts) >= 3 {
		id64, err := strconv.ParseUint(pathParts[2], 10, 32)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid user ID format")
			return
		}
		id := uint(id64)

		if len(pathParts) == 3 {
			switch r.Method {
			case http.MethodGet:
				h.GetUser(w, r, id)
			case http.MethodPut:
				h.UpdateUser(w, r, id)
			case http.MethodDelete:
				h.DeactivateUser(w, r, id)
			default:
				respondError(w, http.StatusMethodNotAllowed, "Method not allowed for /api/users/{id}")
			}
			return
		} else if len(pathParts) == 4 {
			action := pathParts[3]
			switch action {
			case "activate":
				if r.Method == http.MethodPut || r.Method == http.MethodPost {
					h.ActivateUser(w, r, id)
				} else {
					respondError(w, http.StatusMethodNotAllowed, "Use PUT or POST for /api/users/{id}/activate")
				}
			case "permissions":
				switch r.Method {
				case http.MethodGet:
					h.GetUserPermissions(w, r, id)
				case http.MethodPut:
					h.AssignPermissionsToUser(w, r, id)
				default:
					respondError(w, http.StatusMethodNotAllowed, "Use GET or PUT for /api/users/{id}/permissions")
				}
			default:
				http.NotFound(w, r)
			}
			return
		}
	}
	http.NotFound(w, r)
}
