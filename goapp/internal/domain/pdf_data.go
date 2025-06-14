package domain

// Using existing domain.Sale for SaleReceiptData.
// If time.Time is needed for formatting directly in these DTOs, import "time"

type BusinessConfigDetails struct {
	Name     string `json:"nombre" gorm:"column:nombre_empresa"` // Assuming column names from configuracion table
	Phone    string `json:"telefono" gorm:"column:telefono"`
	Address  string `json:"direccion" gorm:"column:direccion"`
	Email    string `json:"email" gorm:"column:email_empresa"`
	LogoPath string `json:"logo_path" gorm:"column:logo"` // Path to logo, relative to a known dir or absolute
	TaxID    string `json:"cuit_ruc" gorm:"column:cuit_ruc"` // For CUIT/RUC
	City     string `json:"ciudad" gorm:"column:ciudad"`
	Country  string `json:"pais" gorm:"column:pais"`
}

type EyePrescriptionPDFDetails struct { // From 'graduaciones' table
	ID           uint   `json:"id_graduacion,omitempty" gorm:"column:id_graduacion;primaryKey"`
	SaleID       uint   `json:"id_venta,omitempty" gorm:"column:id_venta;uniqueIndex"` // Link to the sale
	AddG         string `json:"addg,omitempty" gorm:"column:addg"`
	OD_L1        string `json:"od_l_1,omitempty" gorm:"column:od_l_1"` // Ojo Derecho Lejos - Esferico
	OD_L2        string `json:"od_l_2,omitempty" gorm:"column:od_l_2"` // Ojo Derecho Lejos - Cilindrico
	OD_L3        string `json:"od_l_3,omitempty" gorm:"column:od_l_3"` // Ojo Derecho Lejos - Eje
	OD_C1        string `json:"od_c_1,omitempty" gorm:"column:od_c_1"` // Ojo Derecho Cerca - Esferico
	OD_C2        string `json:"od_c_2,omitempty" gorm:"column:od_c_2"` // Ojo Derecho Cerca - Cilindrico
	OD_C3        string `json:"od_c_3,omitempty" gorm:"column:od_c_3"` // Ojo Derecho Cerca - Eje
	OI_L1        string `json:"oi_l_1,omitempty" gorm:"column:oi_l_1"` // Ojo Izquierdo Lejos - Esferico
	OI_L2        string `json:"oi_l_2,omitempty" gorm:"column:oi_l_2"` // Ojo Izquierdo Lejos - Cilindrico
	OI_L3        string `json:"oi_l_3,omitempty" gorm:"column:oi_l_3"` // Ojo Izquierdo Lejos - Eje
	OI_C1        string `json:"oi_c_1,omitempty" gorm:"column:oi_c_1"` // Ojo Izquierdo Cerca - Esferico
	OI_C2        string `json:"oi_c_2,omitempty" gorm:"column:oi_c_2"` // Ojo Izquierdo Cerca - Cilindrico
	OI_C3        string `json:"oi_c_3,omitempty" gorm:"column:oi_c_3"` // Ojo Izquierdo Cerca - Eje
	Observations string `json:"observaciones,omitempty" gorm:"column:observaciones"`
	// CreatedAt    time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"` // Optional
	// UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"` // Optional
}

// SaleReceiptData aggregates all data needed for a sales receipt PDF.
type SaleReceiptData struct {
	Config       *BusinessConfigDetails
	Sale         *Sale // domain.Sale, expected to be fully populated (Client, User, SaleItems.Product, Payments)
	Prescription *EyePrescriptionPDFDetails // Nullable
}
