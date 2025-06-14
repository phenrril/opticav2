package domain

import (
	"errors"
	"time"
)

var ErrSaleNotFound = errors.New("sale not found")
var ErrInsufficientStock = errors.New("insufficient stock for product")
var ErrPaymentProcessingFailed = errors.New("payment processing failed")
var ErrInvalidSaleData = errors.New("invalid sale data provided")

type Sale struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	ClientID         int       `json:"id_cliente" gorm:"column:id_cliente"`
	UserID           int       `json:"id_usuario" gorm:"column:id_usuario"`
	SaleDate         time.Time `json:"fecha" gorm:"column:fecha"`
	TotalAmount      float64   `json:"total_venta" gorm:"column:total_venta"`
	DiscountAmount   float64   `json:"descuento_monto" gorm:"column:descuento_monto"`
	FinalAmount      float64   `json:"monto_final" gorm:"column:monto_final"`
	AmountPaid       float64   `json:"monto_abonado" gorm:"column:monto_abonado"`
	BalanceDue       float64   `json:"saldo_restante" gorm:"column:saldo_restante"`
	PaymentMethodID  int       `json:"id_metodo_pago" gorm:"column:id_metodo_pago"`
	Status           string    `json:"estado_venta" gorm:"column:estado_venta"`
	PrescriptionData string    `json:"datos_graduacion,omitempty" gorm:"type:json;column:datos_graduacion"`
	Observations     string    `json:"observaciones,omitempty" gorm:"column:observaciones"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	SaleItems []SaleItem `json:"items" gorm:"foreignKey:SaleID"`
	Payments  []Payment  `json:"payments,omitempty" gorm:"foreignKey:SaleID"`
	Client    *Client    `json:"cliente,omitempty" gorm:"foreignKey:ClientID"`
	User      *User      `json:"usuario,omitempty" gorm:"foreignKey:UserID"`
}

type SaleItem struct {
	ID                 int      `json:"id" gorm:"primaryKey"`
	SaleID             int      `json:"id_venta" gorm:"column:id_venta"`
	ProductID          int      `json:"id_producto" gorm:"column:id_producto"`
	Quantity           int      `json:"cantidad" gorm:"column:cantidad"`
	UnitPrice          float64  `json:"precio_unitario" gorm:"column:precio_unitario"`
	TotalPrice         float64  `json:"precio_total_item" gorm:"column:precio_total_item"`
	ProductDescription string   `json:"producto_descripcion,omitempty" gorm:"column:producto_descripcion"`
	ProductCode        string   `json:"producto_codigo,omitempty" gorm:"column:producto_codigo"`
	Product            *Product `json:"producto,omitempty" gorm:"foreignKey:ProductID"`
}

type Payment struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	SaleID            int       `json:"id_venta" gorm:"column:id_venta"`
	PaymentDate       time.Time `json:"fecha_pago" gorm:"column:fecha_pago"`
	Amount            float64   `json:"monto_pago" gorm:"column:monto_pago"`
	PaymentMethodID   int       `json:"id_metodo_pago" gorm:"column:id_metodo_pago"`
	ProcessedByUserID int       `json:"id_usuario_procesa" gorm:"column:id_usuario_procesa"`
	CreatedAt         time.Time `json:"created_at"`
	User              *User     `json:"procesado_por,omitempty" gorm:"foreignKey:ProcessedByUserID"`
}

type CreateSaleRequest struct {
	ClientID         int                     `json:"id_cliente" binding:"required"`
	PaymentMethodID  int                     `json:"id_metodo_pago" binding:"required"`
	DiscountAmount   float64                 `json:"descuento_monto,omitempty"`
	InitialPayment   float64                 `json:"pago_inicial" binding:"required,gte=0"`
	PrescriptionData string                  `json:"datos_graduacion,omitempty"`
	Observations     string                  `json:"observaciones,omitempty"`
	SaleDate         string                  `json:"fecha_venta,omitempty"`
	Items            []CreateSaleItemRequest `json:"items" binding:"required,dive"`
}

type CreateSaleItemRequest struct {
	ProductID int     `json:"id_producto" binding:"required"`
	Quantity  int     `json:"cantidad" binding:"required,gt=0"`
	UnitPrice float64 `json:"precio_unitario,omitempty"`
}

type AddPaymentRequest struct {
	SaleID          int     `json:"id_venta" binding:"required"`
	Amount          float64 `json:"monto_pago" binding:"required,gt=0"`
	PaymentMethodID int     `json:"id_metodo_pago" binding:"required"`
	PaymentDate     string  `json:"fecha_pago,omitempty"`
}

type SaleRepository interface {
	Create(sale *Sale, items []SaleItem, initialPayment *Payment) error
	GetByID(id int) (*Sale, error)
	GetAll(filters map[string]interface{}) ([]Sale, error)
	Update(sale *Sale) error
	GetTopSellingProducts(limit int, fromDate, toDate *time.Time) ([]map[string]interface{}, error)
	Count() (int64, error)
}

type PaymentRepository interface {
	Create(payment *Payment) error
	GetBySaleID(saleID int) ([]Payment, error)
}
