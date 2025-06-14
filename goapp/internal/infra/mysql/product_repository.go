package mysql

import (
	"errors"

	"opticav2/internal/domain"

	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	return r.DB.Table("producto").Create(product).Error
}

func (r *ProductRepository) FindByCode(code string) (*domain.Product, error) {
	var product domain.Product
	err := r.DB.Table("producto").Where("codigo = ?", code).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetByID(id int) (*domain.Product, error) {
	var product domain.Product
	err := r.DB.Table("producto").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.DB.Table("producto").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
	return r.DB.Table("producto").Save(product).Error
}

func (r *ProductRepository) GetLowStockProducts(threshold int, limit int) ([]domain.Product, error) {
	var products []domain.Product
	err := r.DB.Table("producto").
		Where("existencia <= ? AND estado = 1", threshold).
		Order("existencia ASC").
		Limit(limit).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Product{}).Table("producto").Where("estado = 1").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
