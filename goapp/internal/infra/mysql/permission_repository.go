package mysql

import (
	"errors"
	"opticav2/internal/domain"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domain.PermissionRepository {
	return &PermissionRepository{DB: db}
}

// GetAll retrieves all permissions.
func (r *PermissionRepository) GetAll() ([]domain.Permission, error) {
	var permissions []domain.Permission
	// The table name for permissions is 'permisos' based on rol.php
	if err := r.DB.Table("permisos").Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetByID retrieves a permission by its ID.
// Returns domain.ErrPermissionNotFound if no record is found.
func (r *PermissionRepository) GetByID(id uint) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.DB.Table("permisos").First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPermissionNotFound
		}
		return nil, err
	}
	return &permission, nil
}

// GetByIDs retrieves multiple permissions by their IDs.
// This is useful for validating a list of permission IDs.
// It does not return an error if some IDs are not found, just returns the ones that are.
// If an error occurs during query, it's returned.
func (r *PermissionRepository) GetByIDs(ids []uint) ([]domain.Permission, error) {
	var permissions []domain.Permission
	if len(ids) == 0 {
		return permissions, nil // Return empty slice if no IDs provided
	}
	if err := r.DB.Table("permisos").Where("id IN ?", ids).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
