package domain

import (
	"errors"
	"time"
)

// Custom errors for sales domain
var ErrSaleNotFound = errors.New("sale not found")
var ErrInsufficientStock = errors.New("insufficient stock for product")
var ErrPaymentProcessingFailed = errors.New("payment processing failed")
var ErrInvalidSaleData = errors.New("invalid sale data provided")

type Sale struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	ClientID         uint      `json:"id_cliente" gorm:"column:id_cliente"`
	UserID           uint      `json:"id_usuario" gorm:"column:id_usuario"` // User who made the sale
	SaleDate         time.Time `json:"fecha" gorm:"column:fecha"`
	TotalAmount      float64   `json:"total_venta" gorm:"column:total_venta"`     // Sum of SaleItem.TotalPrice
	DiscountAmount   float64   `json:"descuento_monto" gorm:"column:descuento_monto"`
	FinalAmount      float64   `json:"monto_final" gorm:"column:monto_final"`       // TotalAmount - DiscountAmount
	AmountPaid       float64   `json:"monto_abonado" gorm:"column:monto_abonado"`   // Sum of all payments for this sale
	BalanceDue       float64   `json:"saldo_restante" gorm:"column:saldo_restante"` // FinalAmount - AmountPaid
	PaymentMethodID  uint      `json:"id_metodo_pago" gorm:"column:id_metodo_pago"` // Initial payment method
	Status           string    `json:"estado_venta" gorm:"column:estado_venta"`     // e.g., "Pending", "Completed", "Partial", "Cancelled"
	PrescriptionData string    `json:"datos_graduacion,omitempty" gorm:"type:json;column:datos_graduacion"`
	Observations     string    `json:"observaciones,omitempty" gorm:"column:observaciones"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	SaleItems []SaleItem `json:"items" gorm:"foreignKey:SaleID"`
	Payments  []Payment  `json:"payments,omitempty" gorm:"foreignKey:SaleID"` // omitempty as payments might be fetched separately
	Client    *Client    `json:"cliente,omitempty" gorm:"foreignKey:ClientID"` // Optional: for eager loading client details
	User      *User      `json:"usuario,omitempty" gorm:"foreignKey:UserID"`   // Optional: for eager loading user details
}

type SaleItem struct {
	ID                 uint    `json:"id" gorm:"primaryKey"`
	SaleID             uint    `json:"id_venta" gorm:"column:id_venta"`
	ProductID          uint    `json:"id_producto" gorm:"column:id_producto"`
	Quantity           int     `json:"cantidad" gorm:"column:cantidad"`
	UnitPrice          float64 `json:"precio_unitario" gorm:"column:precio_unitario"` // Price at time of sale
	TotalPrice         float64 `json:"precio_total_item" gorm:"column:precio_total_item"`
	ProductDescription string  `json:"producto_descripcion,omitempty" gorm:"column:producto_descripcion"` // Denormalized
	ProductCode        string  `json:"producto_codigo,omitempty" gorm:"column:producto_codigo"`             // Denormalized
	Product            *Product `json:"producto,omitempty" gorm:"foreignKey:ProductID"` // Optional: for eager loading product details
}

type Payment struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	SaleID            uint      `json:"id_venta" gorm:"column:id_venta"`
	PaymentDate       time.Time `json:"fecha_pago" gorm:"column:fecha_pago"`
	Amount            float64   `json:"monto_pago" gorm:"column:monto_pago"`
	PaymentMethodID   uint      `json:"id_metodo_pago" gorm:"column:id_metodo_pago"` // Link to a PaymentMethods table/enum
	ProcessedByUserID uint      `json:"id_usuario_procesa" gorm:"column:id_usuario_procesa"` // User who processed this specific payment
	CreatedAt         time.Time `json:"created_at"`
	User              *User     `json:"procesado_por,omitempty" gorm:"foreignKey:ProcessedByUserID"` // Optional
}

// --- Request Structs ---

type CreateSaleRequest struct {
	ClientID         uint                      `json:"id_cliente" binding:"required"`
	PaymentMethodID  uint                      `json:"id_metodo_pago" binding:"required"`
	DiscountAmount   float64                   `json:"descuento_monto,omitempty"` // Optional discount
	InitialPayment   float64                   `json:"pago_inicial" binding:"required,gte=0"`
	PrescriptionData string                    `json:"datos_graduacion,omitempty"`
	Observations     string                    `json:"observaciones,omitempty"`
	SaleDate         string                    `json:"fecha_venta,omitempty"` // e.g. "YYYY-MM-DD HH:MM:SS", if not provided, use current time
	Items            []CreateSaleItemRequest   `json:"items" binding:"required,dive"`
}

type CreateSaleItemRequest struct {
	ProductID uint    `json:"id_producto" binding:"required"`
	Quantity  int     `json:"cantidad" binding:"required,gt=0"`
	UnitPrice float64 `json:"precio_unitario,omitempty"` // If blank, fetch from product master. If provided, use this price.
}

type AddPaymentRequest struct {
	SaleID            uint    `json:"id_venta" binding:"required"` // To associate payment with a sale
	Amount            float64 `json:"monto_pago" binding:"required,gt=0"`
	PaymentMethodID   uint    `json:"id_metodo_pago" binding:"required"`
	PaymentDate       string  `json:"fecha_pago,omitempty"` // e.g., "YYYY-MM-DD HH:MM:SS", if not provided, use current time
}

// --- Repository Interfaces ---

type SaleRepository interface {
	Create(sale *Sale, items []SaleItem, initialPayment *Payment) error // Transactional
	GetByID(id uint) (*Sale, error)
	// GetAll can be complex. For now, a simple version.
	// Filters could include ClientID, UserID, date range, status.
	GetAll(filters map[string]interface{}) ([]Sale, error)
	Update(sale *Sale) error // For status, observations, etc. Payments handled by PaymentRepository.
	GetTopSellingProducts(limit int, fromDate, toDate *time.Time) ([]map[string]interface{}, error)
	Count() (int64, error)
	// UpdateStatus(id uint, status string) error
	// AddItems(saleID uint, items []SaleItem) error // If items can be added later
	// RemoveItems(saleID uint, itemIDs []uint) error // If items can be removed
}

type PaymentRepository interface {
	Create(payment *Payment) error
	GetBySaleID(saleID uint) ([]Payment, error)
	// GetByID(id uint) (*Payment, error)
}
