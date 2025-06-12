package application

import (
	"opticav2/internal/domain"
	// "encoding/json" // For alternative if prescription is in Sale.PrescriptionData
)

// PrescriptionService defines the interface for accessing prescription data.
type PrescriptionService interface {
	GetPrescriptionForSale(saleID uint) (*domain.EyePrescriptionPDFDetails, error)
	// GetPrescriptionFromJSON(jsonData string) (*domain.EyePrescriptionPDFDetails, error) // Alternative
}

// PrescriptionServiceImpl is the concrete implementation of PrescriptionService.
type PrescriptionServiceImpl struct {
	PrescriptionRepo domain.PrescriptionRepository
	// SaleRepo domain.SaleRepository // If needing to fetch Sale.PrescriptionData JSON
}

// NewPrescriptionService creates a new instance of PrescriptionServiceImpl.
func NewPrescriptionService(prescriptionRepo domain.PrescriptionRepository /*, saleRepo domain.SaleRepository*/) PrescriptionService {
	return &PrescriptionServiceImpl{
		PrescriptionRepo: prescriptionRepo,
		// SaleRepo: saleRepo,
	}
}

// GetPrescriptionForSale retrieves prescription details for a given sale ID.
// This implementation assumes a dedicated 'graduaciones' table.
func (s *PrescriptionServiceImpl) GetPrescriptionForSale(saleID uint) (*domain.EyePrescriptionPDFDetails, error) {
	prescription, err := s.PrescriptionRepo.GetBySaleID(saleID)
	if err != nil {
		// The repository returns nil, nil if not found, so actual errors are DB issues.
		return nil, err
	}
	return prescription, nil // This will be nil if no prescription found, which is acceptable.
}

/*
// Alternative: GetPrescriptionFromJSON parses prescription data from a JSON string.
// This would be used if Sale.PrescriptionData (JSON field) is the source.
func (s *PrescriptionServiceImpl) GetPrescriptionFromJSON(jsonData string) (*domain.EyePrescriptionPDFDetails, error) {
    if jsonData == "" {
        return nil, nil // No JSON data provided
    }
    var prescription domain.EyePrescriptionPDFDetails
    err := json.Unmarshal([]byte(jsonData), &prescription)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling prescription JSON: %w", err)
    }
    return &prescription, nil
}

// Example usage if fetching from Sale.PrescriptionData:
// func (s *PrescriptionServiceImpl) GetPrescriptionForSale(saleID uint) (*domain.EyePrescriptionPDFDetails, error) {
//     sale, err := s.SaleRepo.GetByID(saleID) // Assuming SaleRepo is injected
//     if err != nil {
//         return nil, err
//     }
//     if sale == nil {
//         return nil, domain.ErrSaleNotFound
//     }
//     return s.GetPrescriptionFromJSON(sale.PrescriptionData)
// }
*/
