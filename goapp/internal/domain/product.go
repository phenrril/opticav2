package domain

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrProductCodeTaken = errors.New("product code already exists")

type Product struct {
	ID          int     `json:"id" gorm:"column:codproducto;primaryKey"`
	Code        string  `json:"codigo" gorm:"column:codigo;uniqueIndex"`
	Description string  `json:"descripcion" gorm:"column:descripcion"`
	Brand       string  `json:"marca" gorm:"column:marca"`
	Price       float64 `json:"precio" gorm:"column:precio"`
	Stock       int     `json:"existencia" gorm:"column:existencia"`
	UserID      int     `json:"-" gorm:"column:usuario_id"`
	GrossPrice  float64 `json:"precio_bruto" gorm:"column:precio_bruto"`
	Status      int     `json:"estado" gorm:"column:estado"`
}

type ProductCreateRequest struct {
	Code        string  `json:"codigo" binding:"required"`
	Description string  `json:"descripcion" binding:"required"`
	Brand       string  `json:"marca"`
	Price       float64 `json:"precio" binding:"required"`
	Stock       int     `json:"existencia" binding:"required,gte=0"`
	GrossPrice  float64 `json:"precio_bruto"`
}

type ProductUpdateRequest struct {
	Code        string  `json:"codigo"`
	Description string  `json:"descripcion"`
	Brand       string  `json:"marca"`
	Price       float64 `json:"precio"`
	GrossPrice  float64 `json:"precio_bruto"`
}

type ProductStockUpdateRequest struct {
	AddStock   int     `json:"add_stock"`
	Price      float64 `json:"precio,omitempty"`
	GrossPrice float64 `json:"precio_bruto,omitempty"`
}

type ProductRepository interface {
	Create(product *Product) error
	FindByCode(code string) (*Product, error)
	GetByID(id int) (*Product, error)
	GetAll() ([]Product, error)
	Update(product *Product) error
	GetLowStockProducts(threshold int, limit int) ([]Product, error)
	Count() (int64, error)
}
