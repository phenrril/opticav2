package domain

import "errors"

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
	UserID     int    `json:"-" gorm:"column:usuario_id"`
	Status     int    `json:"estado" gorm:"column:estado"`
	HC         string `json:"hc" gorm:"column:HC"`
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
	Status     *int   `json:"estado"`
	HC         string `json:"hc"`
}

type ClientRepository interface {
	Create(client *Client) error
	FindByName(name string) (*Client, error)
	FindByDNI(dni string) (*Client, error)
	GetByID(id int) (*Client, error)
	GetAll() ([]Client, error)
	Update(client *Client) error
	Count() (int64, error)
}
