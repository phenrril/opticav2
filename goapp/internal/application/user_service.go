package application

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"opticav2/internal/domain"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo       domain.UserRepository
	PermissionRepo domain.PermissionRepository
}

func NewUserService(userRepo domain.UserRepository, permissionRepo domain.PermissionRepository) *UserService {
	return &UserService{UserRepo: userRepo, PermissionRepo: permissionRepo}
}

func (s *UserService) md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *UserService) CreateUser(req domain.UserCreateRequest) (*domain.User, error) {
	_, err := s.UserRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, domain.ErrRecordNotFound) && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("error checking email: " + err.Error())
	}

	_, err = s.UserRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, domain.ErrRecordNotFound) && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("error checking username: " + err.Error())
	}

	hashedPassword := s.md5Hash(req.Password)
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Status:   1,
	}
	errCreate := s.UserRepo.Create(user)
	if errCreate != nil {
		return nil, errCreate
	}
	return user, nil
}

func (s *UserService) GetUser(id uint) (*domain.User, error) {
	return s.UserRepo.GetByID(id)
}

func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.UserRepo.GetAll()
}

func (s *UserService) UpdateUser(id uint, req domain.UserUpdateRequest) (*domain.User, error) {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found for update")
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, errFindByEmail := s.UserRepo.FindByEmail(req.Email)
		if errFindByEmail == nil && existingUser.ID != id {
			return nil, errors.New("email already in use by another account")
		} else if errFindByEmail != nil && !errors.Is(errFindByEmail, domain.ErrRecordNotFound) && !errors.Is(errFindByEmail, gorm.ErrRecordNotFound) {
			return nil, errors.New("error checking email for update: " + errFindByEmail.Error())
		}
	}

	if req.Username != "" && req.Username != user.Username {
		existingUser, errFindByUsername := s.UserRepo.FindByUsername(req.Username)
		if errFindByUsername == nil && existingUser.ID != id {
			return nil, errors.New("username already in use by another account")
		} else if errFindByUsername != nil && !errors.Is(errFindByUsername, domain.ErrRecordNotFound) && !errors.Is(errFindByUsername, gorm.ErrRecordNotFound) {
			return nil, errors.New("error checking username for update: " + errFindByUsername.Error())
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Username = req.Username
	if req.Status != nil {
		user.Status = *req.Status
	}

	errUpdate := s.UserRepo.Update(user)
	if errUpdate != nil {
		return nil, errUpdate
	}
	return user, nil
}

func (s *UserService) DeactivateUser(id uint) error {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found for deactivation")
	}
	user.Status = 0
	return s.UserRepo.Update(user)
}

func (s *UserService) ActivateUser(id uint) error {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found for activation")
	}
	user.Status = 1
	return s.UserRepo.Update(user)
}

func (s *UserService) GetUserPermissions(userID uint) ([]*domain.Permission, error) {
	_, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found when fetching permissions")
	}
	return s.UserRepo.FindPermissionsForUser(userID)
}

func (s *UserService) AssignPermissionsToUser(userID uint, permissionIDs []uint) error {
	_, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found for assigning permissions")
	}

	var permissionsToAssign []*domain.Permission
	if len(permissionIDs) > 0 {
		fetchedPermissionsDomain, errFetchPerms := s.PermissionRepo.GetByIDs(permissionIDs)
		if errFetchPerms != nil {
			return errors.New("error fetching permissions by IDs: " + errFetchPerms.Error())
		}
		if len(fetchedPermissionsDomain) != len(permissionIDs) {
			return errors.New("one or more permission IDs are invalid or not found")
		}
		for i := range fetchedPermissionsDomain {
			permissionsToAssign = append(permissionsToAssign, &fetchedPermissionsDomain[i])
		}
	}

	return s.UserRepo.SetUserPermissions(userID, permissionsToAssign)
}
