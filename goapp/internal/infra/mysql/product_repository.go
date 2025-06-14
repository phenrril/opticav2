package mysql

import (
	"errors"
	"opticav2/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository.
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &ProductRepository{DB: db}
}

// Create creates a new product record.
// GORM will backfill the ID field of the product struct if the insert is successful.
func (r *ProductRepository) Create(product *domain.Product) error {
	// Explicitly use Table("producto") if struct name doesn't match or for clarity
	return r.DB.Table("producto").Create(product).Error
}

// FindByCode finds a product by its code.
// Returns domain.ErrProductNotFound if no record is found.
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

// GetByID retrieves a product by its primary key (codproducto).
// Returns domain.ErrProductNotFound if no record is found.
func (r *ProductRepository) GetByID(id int) (*domain.Product, error) {
	var product domain.Product
	// GORM uses the primary key field defined in the struct tag for First.
	err := r.DB.Table("producto").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products.
// For production, consider adding pagination and filtering.
func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.DB.Table("producto").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Update saves changes to an existing product.
// GORM's Save method updates all fields if the primary key is provided and > 0.
// It will also update 'UpdatedAt' timestamps if available on the model.
func (r *ProductRepository) Update(product *domain.Product) error {
	return r.DB.Table("producto").Save(product).Error
}
