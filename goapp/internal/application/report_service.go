package application

import (
	"fmt"
	"sort"
	"time"

	"opticav2/internal/domain"
	// "gorm.io/gorm" // May be needed for ReportServiceImpl
)

// ReportService defines the interface for generating reports.
type ReportService interface {
	GenerateSalesReport(fromDate, toDate time.Time, userID int, otherFilters map[string]interface{}) (*domain.FullSalesReport, error)
	// Add other report types if needed, e.g., InventoryReport, ClientActivityReport
}

// ReportServiceImpl is the concrete implementation of ReportService.
type ReportServiceImpl struct {
	SaleRepo          domain.SaleRepository
	ProductRepo       domain.ProductRepository // For product details like gross price
	GeneralLedgerRepo domain.GeneralLedgerRepository
	// PaymentRepo domain.PaymentRepository // Not directly used in this specific report logic from reporte.php
	// DB *gorm.DB
}

// NewReportService creates a new instance of ReportServiceImpl.
func NewReportService(
	saleRepo domain.SaleRepository,
	productRepo domain.ProductRepository,
	generalLedgerRepo domain.GeneralLedgerRepository,
) ReportService {
	return &ReportServiceImpl{
		SaleRepo:          saleRepo,
		ProductRepo:       productRepo,
		GeneralLedgerRepo: generalLedgerRepo,
	}
}

// GenerateSalesReport generates a sales report based on the provided date range and filters.
func (s *ReportServiceImpl) GenerateSalesReport(fromDate, toDate time.Time, userID int, otherFilters map[string]interface{}) (*domain.FullSalesReport, error) {
	var reportItems []domain.SalesReportItem
	var summary domain.SalesReportSummary

	// Prepare filters for SaleRepository.GetAll()
	// The GetAll method in SaleRepository already accepts a map[string]interface{} for filters.
	// We will add date_from and date_to to this map.
	// UserID filter might be applied here or in the repository based on role (e.g. non-admin can only see their sales)
	// For now, let's assume userID passed is for context/logging, actual filtering by userID in repo if needed.
	saleFilters := make(map[string]interface{})
	for k, v := range otherFilters { // Copy other pre-existing filters
		saleFilters[k] = v
	}
	saleFilters["date_from"] = fromDate
	saleFilters["date_to"] = toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second) // Ensure includes whole toDate

	// 1. Fetch Sale Details
	sales, err := s.SaleRepo.GetAll(saleFilters)
	if err != nil {
		return nil, fmt.Errorf("error fetching sales for report: %w", err)
	}

	for _, sale := range sales {
		for _, item := range sale.SaleItems {
			// Ensure product is preloaded or fetch it if necessary for GrossPrice
			// SaleRepo.GetAll should ideally preload SaleItems.Product
			var productGrossPrice float64
			if item.Product != nil {
				productGrossPrice = item.Product.GrossPrice
			} else {
				// Fallback: Fetch product if not preloaded (less efficient)
				// This indicates a need to ensure ProductRepo.GetByID uses uint
				// and SaleItem.ProductID is uint.
				// For this subtask, we assume SaleItems.Product is preloaded by SaleRepo.GetByID used in SaleRepo.GetAll's loop.
				// Or, more simply, if SaleRepo.GetAll preloads SaleItems.Product
				prod, pErr := s.ProductRepo.GetByID(item.ProductID)
				if pErr == nil {
					productGrossPrice = prod.GrossPrice
				}
				// else: handle error or assume 0 if product details missing, though this implies data integrity issue
			}

			reportItem := domain.SalesReportItem{
				UserID:            sale.UserID,
				SaleID:            sale.ID,
				ProductID:         item.ProductID,
				ProductName:       item.ProductDescription,
				Quantity:          item.Quantity,
				TransactionDate:   sale.SaleDate.Format("2006-01-02"),
				TransactionType:   "VentaItem",
				ProductGrossPrice: productGrossPrice,
				ProductNetPrice:   item.UnitPrice,  // UnitPrice from SaleItem is the net price at time of sale
				SaleItemTotal:     item.TotalPrice, // UnitPrice * Quantity
				OriginalSaleTotal: sale.TotalAmount,
				SaleDiscount:      sale.DiscountAmount,
				SaleFinalAmount:   sale.FinalAmount,
			}
			reportItems = append(reportItems, reportItem)

			summary.TotalGrossSalesValue += productGrossPrice * float64(item.Quantity)
			summary.TotalNetSalesValue += item.TotalPrice
			summary.TotalDiscountsGiven += sale.DiscountAmount // This will sum discount for each item, needs to be sum of unique sale discounts
		}
	}
	// Correcting TotalDiscountsGiven: Sum unique sale discounts
	summary.TotalDiscountsGiven = 0 // Reset
	processedSaleDiscounts := make(map[uint]bool)
	for _, sale := range sales {
		if !processedSaleDiscounts[uint(sale.ID)] {
			summary.TotalDiscountsGiven += sale.DiscountAmount
			processedSaleDiscounts[uint(sale.ID)] = true
		}
	}

	// 2. Fetch General Income (ingresos)
	generalIncomeEntries, err := s.GeneralLedgerRepo.GetEntriesByDateRange("ingreso", fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error fetching general income: %w", err)
	}
	for _, entry := range generalIncomeEntries {
		reportItems = append(reportItems, domain.SalesReportItem{
			UserID:          entry.UserID, // Assuming UserID exists on GeneralLedgerEntry
			TransactionDate: entry.Date.Format("2006-01-02"),
			TransactionType: entry.Type, // Should be "IngresoGeneral"
			Amount:          entry.Amount,
			// Populate other fields like SaleID if entry.SaleID exists and is relevant
		})
	}
	summary.TotalIncome, err = s.GeneralLedgerRepo.GetSumByDateRange("ingreso", fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error summing general income: %w", err)
	}
	// Note: reporte.php adds sale payments to total_ingresos. Here, TotalIncome is from 'ingresos' table.
	// If sales payments should be included, that logic needs to be added (e.g., sum payments from PaymentRepo or from Sale.AmountPaid).
	// For this subtask, matching reporte.php's direct query on 'ingresos' table:
	// $query4=mysqli_query($con,"select sum(ingresos) as total_ingresos from ingresos where fecha between '$from_date' and '$to_date'");
	// This is covered by GeneralLedgerRepo.GetSumByDateRange("ingreso", ...)

	// 3. Fetch General Expenses (egresos)
	generalExpenseEntries, err := s.GeneralLedgerRepo.GetEntriesByDateRange("egreso", fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error fetching general expenses: %w", err)
	}
	for _, entry := range generalExpenseEntries {
		reportItems = append(reportItems, domain.SalesReportItem{
			UserID:          entry.UserID,
			TransactionDate: entry.Date.Format("2006-01-02"),
			TransactionType: entry.Type, // Should be "EgresoGeneral"
			Amount:          entry.Amount,
		})
	}
	summary.TotalExpenses, err = s.GeneralLedgerRepo.GetSumByDateRange("egreso", fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("error summing general expenses: %w", err)
	}

	// 4. Final Calculations for Summary
	summary.CalculatedProfit = summary.TotalNetSalesValue - summary.TotalGrossSalesValue
	summary.NetBalance = summary.TotalIncome - summary.TotalExpenses // Based on general ledger, not sales profit

	// Optional: Sort all reportItems by date if they are mixed
	sort.Slice(reportItems, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02", reportItems[i].TransactionDate)
		t2, _ := time.Parse("2006-01-02", reportItems[j].TransactionDate)
		return t1.Before(t2)
	})

	fullReport := &domain.FullSalesReport{
		Details: reportItems,
		Summary: summary,
		Filters: domain.ReportFilters{
			FromDate: fromDate.Format("2006-01-02"),
			ToDate:   toDate.Format("2006-01-02"),
			UserID:   userID, // Reflecting the context user, not necessarily a filter applied to all data
		},
	}

	return fullReport, nil
}
