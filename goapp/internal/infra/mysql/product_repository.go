package mysql

import (
	"opticav2/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func (r ProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	// Assuming 'producto' is the table name for products.
	// And domain.Product struct fields map correctly or have GORM tags.
	// The original query selected columns: codproducto, codigo, descripcion, precio, existencia
	// These should map to ID, Code, Description, Price, Stock in domain.Product
	err := r.DB.Table("producto").Find(&products).Error // Or just r.DB.Find(&products) if table name is `products`
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r ProductRepository) Create(p domain.Product) error {
	// Assuming 'producto' is the table name.
	// The original query inserted into columns: codigo, descripcion, precio, existencia
	// These should map to Code, Description, Price, Stock in domain.Product
	// GORM will automatically handle the primary key (ID / codproducto) if it's auto-incrementing
	err := r.DB.Table("producto").Create(&p).Error // Or just r.DB.Create(&p)
	return err
}
