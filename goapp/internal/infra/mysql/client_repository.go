package mysql

import (
	"opticav2/internal/domain"

	"gorm.io/gorm"
)

type ClientRepository struct {
	DB *gorm.DB
}

func NewClientRepository(db *gorm.DB) domain.ClientRepository {
	return &ClientRepository{DB: db}
}

func (r *ClientRepository) Create(client *domain.Client) error {
	return r.DB.Create(client).Error
}

func (r *ClientRepository) FindByName(name string) (*domain.Client, error) {
	var client domain.Client
	err := r.DB.Where("nombre = ?", name).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) FindByDNI(dni string) (*domain.Client, error) {
	var client domain.Client
	err := r.DB.Where("dni = ?", dni).First(&client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) GetByID(id int) (*domain.Client, error) {
	var client domain.Client
	err := r.DB.First(&client, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) GetAll() ([]domain.Client, error) {
	var clients []domain.Client
	err := r.DB.Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ClientRepository) Update(client *domain.Client) error {
	return r.DB.Save(client).Error
}

func (r *ClientRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Client{}).Where("estado = 1").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
