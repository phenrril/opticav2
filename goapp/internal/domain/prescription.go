package domain

// PrescriptionRepository defines the interface for fetching eye prescription data.
type PrescriptionRepository interface {
	// GetBySaleID fetches the prescription details linked to a specific sale ID.
	// Assumes the 'graduaciones' table has an 'id_venta' column.
	GetBySaleID(saleID int) (*EyePrescriptionPDFDetails, error)
}
