package domain

import "errors"

// Custom errors for product domain
var ErrProductNotFound = errors.New("product not found")
var ErrProductCodeTaken = errors.New("product code already exists")

type Product struct {
	ID          uint    `json:"id" gorm:"column:codproducto;primaryKey"` // Changed to uint
	Code        string  `json:"codigo" gorm:"column:codigo;uniqueIndex"`
	Description string  `json:"descripcion" gorm:"column:descripcion"`
	Brand       string  `json:"marca" gorm:"column:marca"`
	Price       float64 `json:"precio" gorm:"column:precio"` // Selling price
	Stock       int     `json:"existencia" gorm:"column:existencia"`
	UserID      uint    `json:"-" gorm:"column:usuario_id"`    // Changed to uint; ID of user who registered/last updated
	GrossPrice  float64 `json:"precio_bruto" gorm:"column:precio_bruto"` // Cost price
	Status      int     `json:"estado" gorm:"column:estado"`     // 0 for inactive, 1 for active
}

type ProductCreateRequest struct {
	Code        string  `json:"codigo" binding:"required"`
	Description string  `json:"descripcion" binding:"required"`
	Brand       string  `json:"marca"`
	Price       float64 `json:"precio" binding:"required"`
	Stock       int     `json:"existencia" binding:"required,gte=0"`
	GrossPrice  float64 `json:"precio_bruto"`
}

type ProductUpdateRequest struct { // For general info, not stock quantity or status
	Code        string  `json:"codigo"` // Allow updating code, must check for uniqueness
	Description string  `json:"descripcion"`
	Brand       string  `json:"marca"`
	Price       float64 `json:"precio"`
	GrossPrice  float64 `json:"precio_bruto"`
}

type ProductStockUpdateRequest struct { // For adjusting stock and optionally price
	AddStock    int     `json:"add_stock"` // Can be negative to reduce stock
	Price       float64 `json:"precio,omitempty"`        // New selling price, if changed during stock update
	GrossPrice  float64 `json:"precio_bruto,omitempty"` // New gross price, if changed
}

type ProductRepository interface {
	Create(product *Product) error
	FindByCode(code string) (*Product, error)
	GetByID(id uint) (*Product, error) // Changed id to uint
	GetAll() ([]Product, error)
	Update(product *Product) error
	GetLowStockProducts(threshold int, limit int) ([]Product, error)
	Count() (int64, error)
	// Delete (soft delete) is an update to Status field, so Update method can cover it.
}
