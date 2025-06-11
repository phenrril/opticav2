package domain

type User struct {
	ID       int    `json:"id" gorm:"column:idusuario;primaryKey"`
	Name     string `json:"nombre" gorm:"column:nombre"`
	Email    string `json:"correo" gorm:"column:correo;uniqueIndex"` // Assuming unique constraint
	Username string `json:"usuario" gorm:"column:usuario;uniqueIndex"` // Assuming unique constraint
	Password string `json:"-" gorm:"column:clave"`
	Status   int    `json:"estado" gorm:"column:estado"`
}

// Request struct for creating a user
type UserCreateRequest struct {
	Name     string `json:"nombre"`
	Email    string `json:"correo"`
	Username string `json:"usuario"`
	Password string `json:"clave"`
}

// Request struct for updating a user (password not included)
type UserUpdateRequest struct {
	Name     string `json:"nombre"`
	Email    string `json:"correo"`
	Username string `json:"usuario"`
	Status   *int   `json:"estado"` // Use pointer to distinguish between 0 and not provided
}

type UserRepository interface {
	GetByCredentials(user, password string) (*User, error) // Existing
	Create(user *User) error                               // New for full user struct
	FindByEmail(email string) (*User, error)               // New
	FindByUsername(username string) (*User, error)         // New
	GetByID(id int) (*User, error)                         // New
	GetAll() ([]User, error)                               // New (or update signature if exists)
	Update(user *User) error                               // New
	// Delete/Deactivate is an update to Status field, so Update method can cover it.
}
