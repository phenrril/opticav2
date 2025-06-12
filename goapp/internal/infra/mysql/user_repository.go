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
// Changed id type from int to uint.
func (r UserRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, id).Error // GORM uses primary key here
	if err != nil {
		return nil, err // Consider returning domain.ErrUserNotFound on gorm.ErrRecordNotFound
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

// FindPermissionsForUser retrieves all permissions associated with a user.
func (r UserRepository) FindPermissionsForUser(userID uint) ([]*domain.Permission, error) {
	var user domain.User
	// Preload the Permissions association. GORM will handle fetching them based on the m2m tag.
	// The table name 'detalle_permisos' is specified in the User struct's GORM tag for Permissions.
	err := r.DB.Preload("Permissions").First(&user, userID).Error
	if err != nil {
		// Handle gorm.ErrRecordNotFound if user itself not found, or return empty slice if user found but no permissions.
		// Preload doesn't error if the association is empty.
		return nil, err
	}
	return user.Permissions, nil
}

// SetUserPermissions updates the permissions for a given user.
// It replaces existing permissions with the new set provided.
func (r UserRepository) SetUserPermissions(userID uint, permissions []*domain.Permission) error {
	var user domain.User
	// First, fetch the user to ensure they exist.
	if err := r.DB.First(&user, userID).Error; err != nil {
		return err // Return error if user not found (e.g., gorm.ErrRecordNotFound)
	}

	// Use GORM's Association().Replace() to manage the many-to-many relationship.
	// This will clear existing associations for the user in 'detalle_permisos' table
	// and insert new ones based on the 'permissions' slice.
	// The 'permissions' slice should contain *domain.Permission objects that already exist in the 'permisos' table.
	// If they don't exist, GORM might try to create them depending on configuration, or it might fail.
	// It's safer if the Permission objects (or at least their IDs) are validated to exist beforehand.
	err := r.DB.Model(&user).Association("Permissions").Replace(permissions)
	if err != nil {
		return err
	}
	return nil
}
