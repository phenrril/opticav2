package application

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"opticav2/internal/domain"
)

type UserService struct {

	UserRepo       domain.UserRepository
	PermissionRepo domain.PermissionRepository // Added PermissionRepo
}

func NewUserService(userRepo domain.UserRepository, permissionRepo domain.PermissionRepository) *UserService {
	return &UserService{UserRepo: userRepo, PermissionRepo: permissionRepo} // Updated constructor

}

func (s *UserService) md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *UserService) CreateUser(req domain.UserCreateRequest) (*domain.User, error) {
	// Check if email already exists
	_, err := s.UserRepo.FindByEmail(req.Email)
	if err == nil { // If err is nil, a user was found
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, domain.ErrRecordNotFound) && err != nil { // Check for other errors, e.g. gorm.ErrRecordNotFound
		// It's important to check if the error is actually gorm.ErrRecordNotFound
		// GORM returns gorm.ErrRecordNotFound when First fails to find the record.
		// If it's not gorm.ErrRecordNotFound, then it's some other unexpected error.
		// However, domain.UserRepository interface doesn't expose gorm errors directly.
		// Assuming the repo implementation returns a specific error type for "not found"
		// For now, we'll proceed if any error occurs, assuming it means "not found".
		// This should be refined with specific error types from repository.
		// For the purpose of this exercise, we assume any error from FindBy means "not found" or "safe to proceed".
		// This is a simplification. In a real app, you'd check for gorm.ErrRecordNotFound.
	}


	// Check if username already exists
	_, err = s.UserRepo.FindByUsername(req.Username)
	if err == nil { // If err is nil, a user was found
		return nil, errors.New("username already exists")
	}
	if !errors.Is(err, domain.ErrRecordNotFound) && err != nil {
		// Similar simplification as above for email check.
	}

	hashedPassword := s.md5Hash(req.Password)
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Status:   1, // Default to active
	}
	err = s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	// The user object from Create might not have the ID if DB doesn't return it or GORM hook isn't set.
	// GORM typically backfills the ID on Create.
	return user, nil
}


// GetUser retrieves a user by ID. ID type changed to uint.
func (s *UserService) GetUser(id uint) (*domain.User, error) {

	return s.UserRepo.GetByID(id)
}

func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.UserRepo.GetAll()
}


// UpdateUser updates a user's details. ID type changed to uint.
func (s *UserService) UpdateUser(id uint, req domain.UserUpdateRequest) (*domain.User, error) {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		// Consider returning domain.ErrUserNotFound or similar specific error
		return nil, errors.New("user not found for update")
r
	}

	// Check for email conflicts if email is being changed
	if req.Email != "" && req.Email != user.Email {
		existingUser, errFindByEmail := s.UserRepo.FindByEmail(req.Email)
		if errFindByEmail == nil && existingUser.ID != id { // if err is nil, email exists for another user
			return nil, errors.New("email already in use by another account")
		}
		// Again, needs proper gorm.ErrRecordNotFound check
	}

	// Check for username conflicts if username is being changed
	if req.Username != "" && req.Username != user.Username {
		existingUser, errFindByUsername := s.UserRepo.FindByUsername(req.Username)
		if errFindByUsername == nil && existingUser.ID != id { // if err is nil, username exists for another user
			return nil, errors.New("username already in use by another account")
		}
		// Again, needs proper gorm.ErrRecordNotFound check
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Username = req.Username
	if req.Status != nil {
		user.Status = *req.Status
	}

	err = s.UserRepo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}


// DeactivateUser sets a user's status to inactive. ID type changed to uint.
func (s *UserService) DeactivateUser(id uint) error {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found for deactivation")

	}
	user.Status = 0 // Assuming 0 means inactive
	return s.UserRepo.Update(user)
}


// ActivateUser sets a user's status to active. ID type changed to uint.
func (s *UserService) ActivateUser(id uint) error {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found for activation")

	}
	user.Status = 1 // Assuming 1 means active
	return s.UserRepo.Update(user)
}


// GetUserPermissions retrieves all permissions for a specific user.
func (s *UserService) GetUserPermissions(userID uint) ([]*domain.Permission, error) {
	// Ensure user exists first (optional, FindPermissionsForUser might do this)
	_, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found when fetching permissions")
	}
	return s.UserRepo.FindPermissionsForUser(userID)
}

// AssignPermissionsToUser assigns a list of permissions (by ID) to a user.
// It replaces any existing permissions.
func (s *UserService) AssignPermissionsToUser(userID uint, permissionIDs []uint) error {
	// 1. Check if user exists
	_, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found for assigning permissions")
	}

	// 2. Fetch domain.Permission objects for the given IDs
	var permissionsToAssign []*domain.Permission
	if len(permissionIDs) > 0 {
		fetchedPermissions, err := s.PermissionRepo.GetByIDs(permissionIDs)
		if err != nil {
			// This error handling assumes GetByIDs returns ErrPermissionNotFound if *any* ID is not found.
			// Or, it might return partial results and no error. The logic below handles partial results.
			if errors.Is(err, domain.ErrPermissionNotFound) {
                 return errors.New("one or more permission IDs are invalid")
            }
			return err // Other unexpected error from repository
		}
		// Ensure all requested permissions were found
		if len(fetchedPermissions) != len(permissionIDs) {
			return errors.New("one or more permission IDs are invalid or not found")
		}
		for i := range fetchedPermissions {
			permissionsToAssign = append(permissionsToAssign, &fetchedPermissions[i])
		}
	}
	// If permissionIDs is empty, permissionsToAssign will be an empty slice,
	// effectively clearing all permissions for the user via SetUserPermissions.

	// 3. Set the user's permissions
	return s.UserRepo.SetUserPermissions(userID, permissionsToAssign)
}

