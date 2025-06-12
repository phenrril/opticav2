package mysql

import (
	"errors"
	"fmt"
	"opticav2/internal/domain"
	"time"
	"gorm.io/gorm"
)

type SaleRepository struct {
	DB         *gorm.DB
	ProductRepo domain.ProductRepository
}

func NewSaleRepository(db *gorm.DB, productRepo domain.ProductRepository) domain.SaleRepository {
	return &SaleRepository{DB: db, ProductRepo: productRepo}
}

func (r *SaleRepository) Create(sale *domain.Sale, items []domain.SaleItem, initialPayment *domain.Payment) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sale).Error; err != nil {
			return fmt.Errorf("failed to create sale: %w", err)
		}

		for i := range items {
			items[i].SaleID = sale.ID
			// Product existence and stock check should ideally happen in service layer before starting transaction,
			// but if done here, it must use 'tx' for product repo calls if ProductRepo is made transaction-aware.
			// For now, ProductRepo.GetByID is not transaction-aware by default.
			// The current approach is to update stock directly with 'tx'.

			// Fetch current stock first to ensure atomicity of check-then-decrement
			var currentProduct domain.Product
			if err := tx.Table("producto").Where("codproducto = ?", items[i].ProductID).First(&currentProduct).Error; err != nil {
			    return fmt.Errorf("product with ID %d not found for stock check: %w", items[i].ProductID, err)
			}

			if currentProduct.Stock < items[i].Quantity {
				return fmt.Errorf("insufficient stock for product ID %d (requested: %d, available: %d): %w",
					items[i].ProductID, items[i].Quantity, currentProduct.Stock, domain.ErrInsufficientStock)
			}

			if err := tx.Model(&domain.Product{}).Where("codproducto = ?", items[i].ProductID).Update("existencia", gorm.Expr("existencia - ?", items[i].Quantity)).Error; err != nil {
				return fmt.Errorf("failed to update stock for product ID %d: %w", items[i].ProductID, err)
			}
			if err := tx.Create(&items[i]).Error; err != nil {
				return fmt.Errorf("failed to create sale item for product ID %d: %w", items[i].ProductID, err)
			}
		}
		if initialPayment != nil && initialPayment.Amount > 0 {
			initialPayment.SaleID = sale.ID
			if err := tx.Create(initialPayment).Error; err != nil {
				return fmt.Errorf("failed to create initial payment for sale ID %d: %w", sale.ID, err)
			}
		}
		return nil
	})
}

func (r *SaleRepository) GetByID(id uint) (*domain.Sale, error) {
	var sale domain.Sale
	err := r.DB.Preload("SaleItems.Product").Preload("Client").Preload("User").Preload("Payments").First(&sale, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSaleNotFound
		}
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) GetAll(filters map[string]interface{}) ([]domain.Sale, error) {
	var sales []domain.Sale
	query := r.DB.Model(&domain.Sale{}).Preload("Client").Preload("User")

	if clientID, ok := filters["client_id"]; ok {
		query = query.Where("id_cliente = ?", clientID)
	}
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("id_usuario = ?", userID)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("estado_venta = ?", status)
	}
	if dateFrom, ok := filters["date_from"]; ok {
		query = query.Where("fecha >= ?", dateFrom)
	}
	if dateTo, ok := filters["date_to"]; ok {
		query = query.Where("fecha <= ?", dateTo)
	}
	query = query.Order("fecha DESC")

	err := query.Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *SaleRepository) Update(sale *domain.Sale) error {
	return r.DB.Save(sale).Error
}

func (r *SaleRepository) GetTopSellingProducts(limit int, fromDate, toDate *time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := r.DB.Table("producto as p").
		Select("p.codproducto, p.descripcion, p.codigo, SUM(di.cantidad) as total_cantidad_vendida, SUM(di.cantidad * di.precio_unitario) as total_revenue").
		Joins("JOIN detalle_venta di ON p.codproducto = di.id_producto")

	if fromDate != nil && toDate != nil {
		query = query.Joins("JOIN ventas v ON di.id_venta = v.id").
			Where("v.fecha BETWEEN ? AND ?", fromDate, toDate)
	}

	err := query.Group("p.codproducto, p.descripcion, p.codigo").
		Order("total_cantidad_vendida DESC").
		Limit(limit).
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *SaleRepository) Count() (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Sale{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
