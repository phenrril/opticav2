package mysql

import (
	"errors"
	"opticav2/internal/domain"

	"gorm.io/gorm"
)

type ConfigRepository struct {
	DB *gorm.DB
}

func NewConfigRepository(db *gorm.DB) domain.ConfigRepository {
	return &ConfigRepository{DB: db}
}

// GetConfig fetches the business configuration from the 'configuracion' table.
// It assumes there is only one row in this table (or fetches the first one).
func (r *ConfigRepository) GetConfig() (*domain.BusinessConfigDetails, error) {
	var config domain.BusinessConfigDetails
	// Assuming the table name is 'configuracion' and it has only one row.
	// The GORM tags in domain.BusinessConfigDetails should map to the column names.
	// Example column names used in domain.BusinessConfigDetails:
	// nombre_empresa, telefono, direccion, email_empresa, logo, cuit_ruc, ciudad, pais
	err := r.DB.Table("configuracion").First(&config).Error // Fetches the first record
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("business configuration not found in database")
		}
		return nil, err
	}
	return &config, nil
}
