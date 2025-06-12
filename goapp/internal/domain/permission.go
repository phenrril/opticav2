package domain

import "errors"

// Custom errors for permission domain
var ErrPermissionNotFound = errors.New("permission not found")

type Permission struct {
	ID   uint   `json:"id" gorm:"column:id;primaryKey"` // Assuming 'id' from rol.php for permisos table
	Name string `json:"nombre" gorm:"column:nombre;uniqueIndex"`
	// Description string `json:"descripcion,omitempty" gorm:"column:descripcion"` // Optional
}

type PermissionRepository interface {
	GetAll() ([]Permission, error)
	GetByID(id uint) (*Permission, error)       // Useful for fetching a single permission
	GetByIDs(ids []uint) ([]Permission, error) // For validating a list of permission IDs
	// Create(permission *Permission) error // Optional: if permissions are managed via this system
	// Update(permission *Permission) error // Optional
	// Delete(id uint) error                // Optional
}
