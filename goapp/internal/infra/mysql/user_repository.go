package mysql

import (
	"opticav2/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// GetByCredentials retrieves a user by username and password (hashed).
func (r UserRepository) GetByCredentials(username, password string) (*domain.User, error) {
	var u domain.User
	// The original query was: "SELECT idusuario, nombre, usuario FROM usuario WHERE usuario=? AND clave=MD5(?) AND estado=1"
	// Assuming 'usuario' is the table name.
	// The domain.User struct now has 'Username' field for 'usuario' column.
	err := r.DB.Table("usuario").Where("usuario = ? AND clave = MD5(?) AND estado = 1", username, password).First(&u).Error
	if err != nil {
		return nil, err // GORM returns gorm.ErrRecordNotFound if no record is found
	}
	return &u, nil
}

// Create creates a new user record in the database.
func (r UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

// FindByEmail retrieves a user by email.
func (r UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("correo = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves a user by username.
func (r UserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("usuario = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by their ID.
// It uses the primary key `idusuario` as defined in the User struct GORM tags.
func (r UserRepository) GetByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, id).Error // GORM uses primary key here
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users.
func (r UserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Update saves changes to an existing user.
// GORM's Save method updates all fields or inserts if the record doesn't exist (if primary key is zero).
// Ensure the user object passed in has the correct ID.
func (r UserRepository) Update(user *domain.User) error {
	return r.DB.Save(user).Error
}
