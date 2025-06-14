package domain

import "time"

// GeneralLedgerEntry represents an entry in either an income or expense table.
// Column names are assumed based on typical usage and reporte.php (e.g., 'fecha', 'descripcion').
// The 'monto' column might be named 'ingresos' in 'ingresos' table and 'egresos' in 'egresos' table.
// The repository implementation will need to handle these specific column names.
type GeneralLedgerEntry struct {
	ID          int       `json:"id" gorm:"primaryKey"`                            // Generic ID, actual column name might be id_ingreso or id_egreso
	Amount      float64   `json:"monto" gorm:"column:monto"`                       // This will map to 'ingresos' or 'egresos' column
	Description string    `json:"descripcion,omitempty" gorm:"column:descripcion"` // Common column name
	Date        time.Time `json:"fecha" gorm:"column:fecha"`                       // Common column name
	Type        string    `json:"tipo_transaccion"`                                // "IngresoGeneral" or "EgresoGeneral", set by service/repo
	UserID      int       `json:"id_usuario,omitempty" gorm:"column:id_usuario"`   // If present in tables
	ClientID    int       `json:"id_cliente,omitempty" gorm:"column:id_cliente"`   // If present in tables
	SaleID      int       `json:"id_venta,omitempty" gorm:"column:id_venta"`       // If present in tables
}

// GeneralLedgerRepository defines methods for accessing general income and expense data.
type GeneralLedgerRepository interface {
	// GetEntriesByDateRange fetches entries. entryType should be "ingreso" or "egreso".
	// The implementation will query the corresponding table.
	GetEntriesByDateRange(entryType string, fromDate, toDate time.Time) ([]GeneralLedgerEntry, error)

	// GetSumByDateRange calculates the sum of amounts for a given entry type.
	// entryType should be "ingreso" or "egreso".
	GetSumByDateRange(entryType string, fromDate, toDate time.Time) (float64, error)
}
