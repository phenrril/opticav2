package application

import (
	"opticav2/internal/domain"
)

type PermissionService struct {
	PermissionRepo domain.PermissionRepository
}

func NewPermissionService(repo domain.PermissionRepository) *PermissionService {
	return &PermissionService{PermissionRepo: repo}
}

// ListAll retrieves all available permissions.
func (s *PermissionService) ListAll() ([]domain.Permission, error) {
	permissions, err := s.PermissionRepo.GetAll()
	if err != nil {
		return nil, err // Or wrap the error if needed
	}
	return permissions, nil
}

// GetPermissionByID retrieves a single permission by its ID.
// This might be useful for validating a permission exists before assigning it.
func (s *PermissionService) GetPermissionByID(id uint) (*domain.Permission, error) {
	permission, err := s.PermissionRepo.GetByID(id)
	if err != nil {
		return nil, err // Handles domain.ErrPermissionNotFound from repo
	}
	return permission, nil
}

// GetPermissionsByIDs retrieves multiple permissions by their IDs.
// Useful for converting a list of IDs from a request to a list of Permission objects.
func (s *PermissionService) GetPermissionsByIDs(ids []uint) ([]domain.Permission, error) {
    if len(ids) == 0 {
        return []domain.Permission{}, nil
    }
	permissions, err := s.PermissionRepo.GetByIDs(ids)
	if err != nil {
		return nil, err
	}
    if len(permissions) != len(ids) {
        // This indicates that some permission IDs were not found.
        // Depending on strictness, this could be an error or just a partial result.
        // For assigning permissions, it should likely be an error.
        return permissions, domain.ErrPermissionNotFound // Or a more specific error like "some permissions not found"
    }
	return permissions, nil
}
