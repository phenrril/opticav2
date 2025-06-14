package domain

type User struct {
	ID          uint          `json:"id" gorm:"column:idusuario;primaryKey"`
	Name        string        `json:"nombre" gorm:"column:nombre"`
	Email       string        `json:"correo" gorm:"column:correo;uniqueIndex"`
	Username    string        `json:"usuario" gorm:"column:usuario;uniqueIndex"`
	Password    string        `json:"-" gorm:"column:clave"`
	Status      int           `json:"estado" gorm:"column:estado"`
	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:detalle_permisos;foreignKey:ID;joinForeignKey:id_usuario;References:ID;joinReferences:id_permiso"`
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
	Status   *int   `json:"estado"`
}

type UserRepository interface {
	GetByCredentials(user, password string) (*User, error)
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	GetByID(id uint) (*User, error)
	GetAll() ([]User, error)
	Update(user *User) error
	FindPermissionsForUser(userID uint) ([]*Permission, error)
	SetUserPermissions(userID uint, permissions []*Permission) error
	Count() (int64, error)
}
