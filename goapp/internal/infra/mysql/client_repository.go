package mysql

import (
	"opticav2/internal/domain"
	"gorm.io/gorm"
)

type ClientRepository struct {
	DB *gorm.DB
}

// NewClientRepository creates a new instance of ClientRepository.
// It's good practice to have a constructor.
func NewClientRepository(db *gorm.DB) domain.ClientRepository {
	return &ClientRepository{DB: db}
}

// Create creates a new client record.
func (r *ClientRepository) Create(client *domain.Client) error {
	return r.DB.Create(client).Error
}

// FindByName finds a client by name.
// GORM's First will return gorm.ErrRecordNotFound if no record is found.
func (r *ClientRepository) FindByName(name string) (*domain.Client, error) {
	var client domain.Client
	err := r.DB.Where("nombre = ?", name).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound // Use domain specific error
		}
		return nil, err
	}
	return &client, nil
}

// FindByDNI finds a client by DNI.
func (r *ClientRepository) FindByDNI(dni string) (*domain.Client, error) {
	var client domain.Client
	err := r.DB.Where("dni = ?", dni).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound // Use domain specific error
		}
		return nil, err
	}
	return &client, nil
}

// GetByID retrieves a client by their ID (primary key).
// Changed id type from int to uint.
func (r *ClientRepository) GetByID(id uint) (*domain.Client, error) {
	var client domain.Client
	err := r.DB.First(&client, id).Error // GORM uses primary key here (idcliente)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound // Use domain specific error
		}
		return nil, err
	}
	return &client, nil
}

// GetAll retrieves all clients.
// Consider adding filtering or pagination for production systems.
func (r *ClientRepository) GetAll() ([]domain.Client, error) {
	var clients []domain.Client
	err := r.DB.Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// Update saves changes to an existing client.
// GORM's Save method updates all fields if the primary key is provided and > 0.
func (r *ClientRepository) Update(client *domain.Client) error {
	return r.DB.Save(client).Error
}

// Count returns the total number of active clients.
func (r *ClientRepository) Count() (int64, error) {
	var count int64
	// Assuming active clients have 'estado = 1'
	err := r.DB.Model(&domain.Client{}).Where("estado = 1").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
