package domain

// ConfigRepository defines the interface for fetching business configuration.
type ConfigRepository interface {
	// GetConfig fetches the single row of business configuration.
	// Assumes there's only one configuration entry in the table.
	GetConfig() (*BusinessConfigDetails, error)
}
