package mysql

import (
	"errors"
	"opticav2/internal/domain"

	"gorm.io/gorm"
)

type PrescriptionRepository struct {
	DB *gorm.DB
}

func NewPrescriptionRepository(db *gorm.DB) domain.PrescriptionRepository {
	return &PrescriptionRepository{DB: db}
}

// GetBySaleID fetches prescription details from the 'graduaciones' table
// based on the associated sale ID.
func (r *PrescriptionRepository) GetBySaleID(saleID uint) (*domain.EyePrescriptionPDFDetails, error) {
	var prescription domain.EyePrescriptionPDFDetails
	// The GORM tags in domain.EyePrescriptionPDFDetails should map to column names in 'graduaciones'.
	// We are querying by 'id_venta'.
	err := r.DB.Table("graduaciones").Where("id_venta = ?", saleID).First(&prescription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// It's okay for a sale not to have a prescription.
			// Return nil, nil to indicate "not found but not an error".
			return nil, nil
		}
		return nil, err // Other actual database error
	}
	return &prescription, nil
}
