package middleware

import (
	"context"
	"net/http"
	"opticav2/internal/application"
	"strconv"
	// "log" // For debugging
)

// UserIDKey is a context key for storing UserID
type contextKey string

const UserIDKey contextKey = "userID"

type AuthorizationMiddleware struct {
	UserService *application.UserService
}

func NewAuthorizationMiddleware(us *application.UserService) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{UserService: us}
}

func (amw *AuthorizationMiddleware) RequirePermission(requiredPermissionName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Identify User (Placeholder: from header)
			userIDStr := r.Header.Get("X-User-ID")
			if userIDStr == "" {
				// log.Println("Authz: Missing X-User-ID header")
				http.Error(w, "Unauthorized: Missing user identifier", http.StatusUnauthorized)
				return
			}
			userID64, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				// log.Printf("Authz: Invalid X-User-ID header: %s\n", userIDStr)
				http.Error(w, "Unauthorized: Invalid user identifier format", http.StatusUnauthorized)
				return
			}
			userID := uint(userID64)

			// Store userID in context for potential use by handlers
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			// 2. Fetch User's Permissions (Superadmin User ID 1 bypasses checks)
			if userID == 1 { // Superadmin bypass
				// log.Printf("Authz: UserID 1 (Superadmin) granted access to %s for permission %s\n", r.URL.Path, requiredPermissionName)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			permissions, err := amw.UserService.GetUserPermissions(userID)
			if err != nil {
				// log.Printf("Authz: Error fetching permissions for UserID %d: %v\n", userID, err)
				// Consider the case where user is not found vs. other errors
				// If UserService.GetUserPermissions returns a specific "user not found" error:
				// if errors.Is(err, domain.ErrUserNotFound) { // Assuming domain.ErrUserNotFound exists
				//    http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
				//    return
				// }
				http.Error(w, "Forbidden: Could not retrieve user permissions", http.StatusForbidden)
				return
			}

			// 3. Check for Required Permission
			hasPermission := false
			for _, p := range permissions {
				if p.Name == requiredPermissionName {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				// log.Printf("Authz: UserID %d denied access to %s, lacks permission %s\n", userID, r.URL.Path, requiredPermissionName)
				http.Error(w, "Forbidden: You don't have the required permission ("+requiredPermissionName+")", http.StatusForbidden)
				return
			}

			// log.Printf("Authz: UserID %d granted access to %s for permission %s\n", userID, r.URL.Path, requiredPermissionName)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
