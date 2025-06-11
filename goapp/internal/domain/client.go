package domain

import "errors" // Added import for errors package

// It's good practice to define error variables for common errors
var ErrClientNotFound = errors.New("client not found")
var ErrClientDNITaken = errors.New("client DNI already exists")
var ErrClientNameTaken = errors.New("client name already exists (DNI not provided or also taken)")


type Client struct {
	ID         int    `json:"id" gorm:"column:idcliente;primaryKey"`
	Name       string `json:"nombre" gorm:"column:nombre"`
	Phone      string `json:"telefono" gorm:"column:telefono"`
	Address    string `json:"direccion" gorm:"column:direccion"`
	DNI        string `json:"dni" gorm:"column:dni;uniqueIndex"`
	ObraSocial string `json:"obrasocial" gorm:"column:obrasocial"`
	Medico     string `json:"medico" gorm:"column:medico"`
	UserID     int    `json:"-" gorm:"column:usuario_id"` // ID of user who registered
	Status     int    `json:"estado" gorm:"column:estado"`
	HC         string `json:"hc" gorm:"column:HC"`
	// User    User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:idusuario"` // Corrected reference if User table PK is idusuario
}

type ClientCreateRequest struct {
	Name       string `json:"nombre" binding:"required"`
	Phone      string `json:"telefono"`
	Address    string `json:"direccion"`
	DNI        string `json:"dni" binding:"required"`
	ObraSocial string `json:"obrasocial"`
	Medico     string `json:"medico"`
	HC         string `json:"hc"`
}

type ClientUpdateRequest struct {
	Name       string `json:"nombre"`
	Phone      string `json:"telefono"`
	Address    string `json:"direccion"`
	DNI        string `json:"dni"`
	ObraSocial string `json:"obrasocial"`
	Medico     string `json:"medico"`
	Status     *int   `json:"estado"` // Pointer for explicit update to distinguish 0 from not provided
	HC         string `json:"hc"`
}

type ClientRepository interface {
	Create(client *Client) error
	FindByName(name string) (*Client, error) // Consider if this should return a list if names aren't unique
	FindByDNI(dni string) (*Client, error)
	GetByID(id int) (*Client, error)
	GetAll() ([]Client, error)
	Update(client *Client) error
	// Delete is typically an update to Status (soft delete)
}
