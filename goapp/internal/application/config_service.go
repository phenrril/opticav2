package application

import (
	"opticav2/internal/domain"
)

// ConfigService defines the interface for accessing business configuration.
type ConfigService interface {
	GetBusinessDetails() (*domain.BusinessConfigDetails, error)
}

// ConfigServiceImpl is the concrete implementation of ConfigService.
type ConfigServiceImpl struct {
	ConfigRepo domain.ConfigRepository
}

// NewConfigService creates a new instance of ConfigServiceImpl.
func NewConfigService(configRepo domain.ConfigRepository) ConfigService {
	return &ConfigServiceImpl{ConfigRepo: configRepo}
}

// GetBusinessDetails retrieves the business configuration details.
func (s *ConfigServiceImpl) GetBusinessDetails() (*domain.BusinessConfigDetails, error) {
	config, err := s.ConfigRepo.GetConfig()
	if err != nil {
		// Potentially wrap error or handle specific cases
		return nil, err
	}
	return config, nil
}
