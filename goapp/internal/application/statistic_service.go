package application

import (
	"fmt"
	"strconv"
	"time"

	"opticav2/internal/domain"
	// "gorm.io/gorm" // May be needed for StatisticServiceImpl
)

// StatisticService defines the interface for fetching various statistics.
type StatisticService interface {
	GetLowStockProducts(threshold int) ([]domain.LowStockProductStat, error)
	GetTopSellingProducts(limit int, fromDateStr, toDateStr string) ([]domain.TopSellingProductStat, error)
	GetDashboardSummaryCounts() (*domain.DashboardSummaryCounts, error)
	// Add other statistic methods as needed
}

// StatisticServiceImpl is the concrete implementation of StatisticService.
// It will require access to various repositories.
type StatisticServiceImpl struct {
	ProductRepo domain.ProductRepository
	SaleRepo    domain.SaleRepository   // For top selling products, sale counts
	UserRepo    domain.UserRepository   // For user counts
	ClientRepo  domain.ClientRepository // For client counts
	// DB *gorm.DB // For complex custom queries if needed
}

// NewStatisticService creates a new instance of StatisticServiceImpl.
func NewStatisticService(
	productRepo domain.ProductRepository,
	saleRepo domain.SaleRepository,
	userRepo domain.UserRepository,
	clientRepo domain.ClientRepository,
// db *gorm.DB,
) StatisticService {
	return &StatisticServiceImpl{
		ProductRepo: productRepo,
		SaleRepo:    saleRepo,
		UserRepo:    userRepo,
		ClientRepo:  clientRepo,
		// DB: db,
	}
}

// GetLowStockProducts retrieves products with stock less than or equal to a given threshold.
func (s *StatisticServiceImpl) GetLowStockProducts(threshold int) ([]domain.LowStockProductStat, error) {
	// Define a reasonable limit for how many low stock products to return, e.g., 50
	limit := 50
	products, err := s.ProductRepo.GetLowStockProducts(threshold, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching low stock products: %w", err)
	}

	stats := make([]domain.LowStockProductStat, len(products))
	for i, p := range products {
		stats[i] = domain.LowStockProductStat{
			ProductID:     p.ID,
			ProductName:   p.Description, // Assuming Description is the main name field
			ProductCode:   p.Code,
			StockQuantity: p.Stock,
			BrandName:     p.Brand,
		}
	}
	return stats, nil
}

// GetTopSellingProducts retrieves products ranked by their sales quantity within a date range.
func (s *StatisticServiceImpl) GetTopSellingProducts(limit int, fromDateStr, toDateStr string) ([]domain.TopSellingProductStat, error) {
	var fromDate, toDate *time.Time

	if fromDateStr != "" {
		t, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid from_date format, use YYYY-MM-DD: %w", err)
		}
		fromDate = &t
	}
	if toDateStr != "" {
		t, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid to_date format, use YYYY-MM-DD: %w", err)
		}
		// Adjust toDate to include the whole day
		t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		toDate = &t
	}

	results, err := s.SaleRepo.GetTopSellingProducts(limit, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error fetching top selling products: %w", err)
	}

	stats := make([]domain.TopSellingProductStat, len(results))
	for i, res := range results {
		stat := domain.TopSellingProductStat{}
		if id, ok := res["codproducto"]; ok {
			if idFloat, ok := id.(float64); ok { // SUM often returns float64 from DB
				stat.ProductID = int(idFloat)
			} else if idInt, ok := id.(int64); ok {
				stat.ProductID = int(idInt)
			}
		}
		if name, ok := res["descripcion"].(string); ok {
			stat.ProductName = name
		}
		if code, ok := res["codigo"].(string); ok {
			stat.ProductCode = code
		}
		var qtyStr []byte
		if qty, ok := res["total_cantidad_vendida"]; ok {
			// SUM usually returns a type that can be float64 or int64 from DB drivers
			if qtyFloat, ok := qty.(float64); ok {
				stat.TotalSoldQuantity = int(qtyFloat)
			} else if qtyInt, ok := qty.(int64); ok {
				stat.TotalSoldQuantity = int(qtyInt)
			} else if tmp, ok := qty.([]byte); ok {
				qtyStr = tmp
				qtyInt, _ := strconv.Atoi(string(qtyStr))
				stat.TotalSoldQuantity = qtyInt
			}
		}
		var revStr []byte
		if revenue, ok := res["total_revenue"]; ok {
			if revFloat, ok := revenue.(float64); ok {
				stat.TotalRevenue = revFloat
			} else if tmp, ok := revenue.([]byte); ok {
				revStr = tmp
				revFloat, _ := strconv.ParseFloat(string(revStr), 64)
				stat.TotalRevenue = revFloat
			}
		}
		stats[i] = stat
	}
	return stats, nil
}

// GetDashboardSummaryCounts retrieves total counts for key entities.
func (s *StatisticServiceImpl) GetDashboardSummaryCounts() (*domain.DashboardSummaryCounts, error) {
	userCount, err := s.UserRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("error counting users: %w", err)
	}
	clientCount, err := s.ClientRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("error counting clients: %w", err)
	}
	productCount, err := s.ProductRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("error counting products: %w", err)
	}
	saleCount, err := s.SaleRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("error counting sales: %w", err)
	}

	return &domain.DashboardSummaryCounts{
		UserCount:    userCount,
		ClientCount:  clientCount,
		ProductCount: productCount,
		SaleCount:    saleCount,
	}, nil
}
