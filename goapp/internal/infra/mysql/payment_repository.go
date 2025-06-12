package mysql

import (
	// "errors"
	"opticav2/internal/domain"
	"gorm.io/gorm"
	// "gorm.io/gorm/clause" // If needed for specific clauses like OnConflict
)

type PaymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &PaymentRepository{DB: db}
}

// Create creates a new payment record.
// This might also need to update the parent Sale's AmountPaid and BalanceDue fields.
// Typically, this logic is handled in the service layer after successful payment creation.
func (r *PaymentRepository) Create(payment *domain.Payment) error {
	return r.DB.Create(payment).Error
}

// GetBySaleID retrieves all payments associated with a given SaleID.
func (r *PaymentRepository) GetBySaleID(saleID uint) ([]domain.Payment, error) {
	var payments []domain.Payment
	// Optionally preload user who processed payment if that's needed in the list.
	err := r.DB.Preload("User").Where("id_venta = ?", saleID).Order("fecha_pago ASC").Find(&payments).Error
	if err != nil {
		// It's okay to return an empty slice if no payments found, not necessarily an error.
		// gorm.ErrRecordNotFound is not returned by Find when the result is a slice.
		return nil, err
	}
	return payments, nil
}

// GetByID could be added if direct access to a payment by its ID is needed.
// func (r *PaymentRepository) GetByID(id uint) (*domain.Payment, error) {
// 	var payment domain.Payment
// 	err := r.DB.First(&payment, id).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("payment not found") // Or a domain specific error
// 		}
// 		return nil, err
// 	}
// 	return &payment, nil
// }
