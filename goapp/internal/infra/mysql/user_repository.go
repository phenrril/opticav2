package mysql

import (
	"opticav2/internal/domain"
	"gorm.io/gorm"
	// "errors"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r UserRepository) GetByCredentials(username, password string) (*domain.User, error) {
	var u domain.User
	err := r.DB.Table("usuario").Where("usuario = ? AND clave = MD5(?) AND estado = 1", username, password).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("correo = ?", email).First(&user).Error
	// Propagate gorm.ErrRecordNotFound for service to check, or convert to domain.ErrUserNotFound
	// if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domain.ErrRecordNotFound } // Assuming domain.ErrRecordNotFound exists
	return &user, err
}

func (r UserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("usuario = ?", username).First(&user).Error
	// Propagate gorm.ErrRecordNotFound for service to check
	// if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domain.ErrRecordNotFound }
	return &user, err
}

func (r UserRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, id).Error
	// if errors.Is(err, gorm.ErrRecordNotFound) { return nil, domain.ErrUserNotFound } // Assuming domain.ErrUserNotFound
	return &user, err
}

func (r UserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r UserRepository) Update(user *domain.User) error {
	return r.DB.Save(user).Error
}

func (r UserRepository) FindPermissionsForUser(userID uint) ([]*domain.Permission, error) {
	var user domain.User
	err := r.DB.Preload("Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Permissions, nil
}

func (r UserRepository) SetUserPermissions(userID uint, permissions []*domain.Permission) error {
	var user domain.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return err
	}
	err := r.DB.Model(&user).Association("Permissions").Replace(permissions)
	return err
}

func (r UserRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&domain.User{}).Where("estado = 1").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
