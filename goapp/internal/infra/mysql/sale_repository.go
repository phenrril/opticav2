package mysql

import (
	"errors"
	"fmt"
	"opticav2/internal/domain"

	"gorm.io/gorm"
)

type SaleRepository struct {
	DB         *gorm.DB
	ProductRepo domain.ProductRepository // Inject ProductRepository for stock updates
}

func NewSaleRepository(db *gorm.DB, productRepo domain.ProductRepository) domain.SaleRepository {
	return &SaleRepository{DB: db, ProductRepo: productRepo}
}

// Create handles the creation of a sale, its items, and an initial payment within a transaction.
// It also updates product stock.
func (r *SaleRepository) Create(sale *domain.Sale, items []domain.SaleItem, initialPayment *domain.Payment) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create the Sale record
		if err := tx.Create(sale).Error; err != nil {
			return fmt.Errorf("failed to create sale: %w", err)
		}

		// 2. Create SaleItem records and update product stock
		for i := range items {
			items[i].SaleID = sale.ID // Link item to the created sale

			// Fetch product for stock update and to confirm existence
			product, err := r.ProductRepo.GetByID(items[i].ProductID) // Use ProductRepo from outside transaction context
			if err != nil {
				return fmt.Errorf("product with ID %d not found for sale item: %w", items[i].ProductID, err)
			}

			if product.Stock < items[i].Quantity {
				return fmt.Errorf("insufficient stock for product ID %d (requested: %d, available: %d): %w",
					items[i].ProductID, items[i].Quantity, product.Stock, domain.ErrInsufficientStock)
			}
			product.Stock -= items[i].Quantity
			// Using r.ProductRepo.Update here will use the original DB instance, not the transaction `tx`.
			// For stock updates to be part of the transaction, ProductRepo methods would need to accept `*gorm.DB`
			// or we perform the update directly with `tx`.
			if err := tx.Model(&domain.Product{}).Where("codproducto = ?", product.ID).Update("existencia", product.Stock).Error; err != nil {
				return fmt.Errorf("failed to update stock for product ID %d: %w", product.ID, err)
			}

			// Denormalize product details if necessary (already part of item struct from service)
			// items[i].ProductDescription = product.Description
			// items[i].ProductCode = product.Code

			if err := tx.Create(&items[i]).Error; err != nil {
				return fmt.Errorf("failed to create sale item for product ID %d: %w", items[i].ProductID, err)
			}
		}

		// Update sale.SaleItems to reflect created items with IDs (GORM might handle this if association is set up)
		// sale.SaleItems = items // This might not be necessary if GORM handles associations correctly on Sale create with nested items.
		// For explicit control, we create items separately as above.

		// 3. Create the initial Payment record, if provided
		if initialPayment != nil && initialPayment.Amount > 0 {
			initialPayment.SaleID = sale.ID // Link payment to the created sale
			if err := tx.Create(initialPayment).Error; err != nil {
				return fmt.Errorf("failed to create initial payment for sale ID %d: %w", sale.ID, err)
			}
		}
		return nil // Commit transaction
	})
}

func (r *SaleRepository) GetByID(id uint) (*domain.Sale, error) {
	var sale domain.Sale
	// Preload SaleItems and associated Product, Client, User for detailed view
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
	query := r.DB.Model(&domain.Sale{}).Preload("Client").Preload("User") // Preload basic info for list views

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
	// Add order by most recent
	query = query.Order("fecha DESC")


	err := query.Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

// Update is for general sale updates like status, observations.
// Payments should be handled via PaymentRepository.
// Item modifications would need more specific methods.
func (r *SaleRepository) Update(sale *domain.Sale) error {
	// Ensure only certain fields are updatable or use DB.Model(&domain.Sale{}).Where("id = ?", sale.ID).Updates(map[string]interface{}{...})
	// For simplicity, Save() updates all fields based on the primary key.
	return r.DB.Save(sale).Error
}

// GetTopSellingProducts retrieves products ranked by their sales quantity.
func (r *SaleRepository) GetTopSellingProducts(limit int, fromDate, toDate *time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Base query selects product details and sums up quantity from sale items
	query := r.DB.Table("producto as p").
		Select("p.codproducto, p.descripcion, p.codigo, SUM(di.cantidad) as total_cantidad_vendida, SUM(di.cantidad * di.precio_unitario) as total_revenue").
		Joins("JOIN detalle_venta di ON p.codproducto = di.id_producto")

	// If date filters are provided, join with sales table and apply date range
	if fromDate != nil && toDate != nil {
		query = query.Joins("JOIN ventas v ON di.id_venta = v.id").
			Where("v.fecha BETWEEN ? AND ?", fromDate, toDate)
	}

	// Group by product information and order by the total quantity sold
	err := query.Group("p.codproducto, p.descripcion, p.codigo").
		Order("total_cantidad_vendida DESC").
		Limit(limit).
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

// Count returns the total number of sales records (optionally filtered by active/completed status).
func (r *SaleRepository) Count() (int64, error) {
	var count int64
	// Example: Count all sales. Could be refined to count "Completed" or "Pending" sales.
	err := r.DB.Model(&domain.Sale{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
