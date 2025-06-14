package domain

// --- Sales Report Structs ---

type SalesReportItem struct { // For individual line items in the report
	UserID            int     `json:"id_usuario,omitempty"`
	SaleID            int     `json:"id_venta,omitempty"` // Might be Sale ID or an Ingreso/Egreso ID
	ProductID         int     `json:"id_producto,omitempty"`
	ProductName       string  `json:"nombre_producto,omitempty"`
	Quantity          int     `json:"cantidad,omitempty"`
	TransactionDate   string  `json:"fecha"`            // Keep as string YYYY-MM-DD for consistency with input
	TransactionType   string  `json:"tipo_transaccion"` // e.g., "VentaItem", "IngresoGeneral", "EgresoGeneral"
	Amount            float64 `json:"monto,omitempty"`  // For Ingreso/Egreso items, or payment amounts
	ProductGrossPrice float64 `json:"precio_bruto_producto,omitempty"`
	ProductNetPrice   float64 `json:"precio_neto_producto,omitempty"`
	SaleItemTotal     float64 `json:"total_item_venta,omitempty"`     // ProductNetPrice * Quantity
	OriginalSaleTotal float64 `json:"total_venta_original,omitempty"` // Total of the parent sale for context (before discount)
	SaleDiscount      float64 `json:"descuento_venta,omitempty"`      // Discount on the parent sale
	SaleFinalAmount   float64 `json:"monto_final_venta,omitempty"`    // Final amount of the parent sale
}

type SalesReportSummary struct {
	TotalIncome          float64 `json:"total_ingresos"`             // Sum of all income (sale payments, IngresoGeneral)
	TotalExpenses        float64 `json:"total_egresos"`              // Sum of all EgresoGeneral
	TotalGrossSalesValue float64 `json:"total_valor_venta_bruta"`    // Sum of (ProductGrossPrice * Quantity) for all sold items
	TotalNetSalesValue   float64 `json:"total_valor_venta_neta"`     // Sum of (ProductNetPrice * Quantity) for all sold items (SaleItemTotal)
	TotalDiscountsGiven  float64 `json:"total_descuentos_otorgados"` // Sum of Sale.DiscountAmount
	CalculatedProfit     float64 `json:"ganancia_calculada"`         // TotalNetSalesValue - TotalGrossSalesValue (Profit from sales only)
	NetBalance           float64 `json:"balance_neto"`               // TotalIncome - TotalExpenses (Overall cash flow)
}

type FullSalesReport struct {
	Details []SalesReportItem  `json:"detalles"`
	Summary SalesReportSummary `json:"resumen"`
	Filters ReportFilters      `json:"filtros_aplicados,omitempty"`
}

type ReportFilters struct {
	FromDate string `json:"fecha_desde,omitempty"`
	ToDate   string `json:"fecha_hasta,omitempty"`
	UserID   int    `json:"id_usuario,omitempty"`
	// Add other filters as needed
}

// --- Statistics API Response Structs ---

type LowStockProductStat struct {
	ProductID     int    `json:"id_producto"`
	ProductName   string `json:"nombre_producto"`
	ProductCode   string `json:"codigo_producto"`
	StockQuantity int    `json:"cantidad_stock"`
	BrandName     string `json:"marca_producto,omitempty"`
}

type TopSellingProductStat struct {
	ProductID         int     `json:"id_producto"`
	ProductName       string  `json:"nombre_producto"`
	ProductCode       string  `json:"codigo_producto"`
	TotalSoldQuantity int     `json:"total_cantidad_vendida"`
	TotalRevenue      float64 `json:"total_ingresos_generados"` // Sum of SaleItem.TotalPrice for this product
}

type DashboardSummaryCounts struct {
	UserCount    int64 `json:"total_usuarios"`  // Changed to int64 for GORM Count
	ClientCount  int64 `json:"total_clientes"`  // Changed to int64
	ProductCount int64 `json:"total_productos"` // Changed to int64
	SaleCount    int64 `json:"total_ventas"`    // Total number of sales transactions
}

// --- Potentially needed for Ingreso/Egreso if they are separate tables ---
// These would be similar to Sale but simpler.
// For now, SalesReportItem.TransactionType and Amount can cover them if queried from a combined source.

// Example:
// type GeneralLedgerEntry struct {
//     ID uint `json:"id" gorm:"primaryKey"`
//     EntryDate time.Time `json:"fecha_entrada" gorm:"column:fecha"`
//     Description string `json:"descripcion" gorm:"column:descripcion"`
//     Amount float64 `json:"monto" gorm:"column:monto"` // Positive for income, negative for expense
//     Type string `json:"tipo" gorm:"column:tipo"` // "Ingreso", "Egreso"
//     UserID uint `json:"id_usuario" gorm:"column:id_usuario"` // User who registered it
//     CreatedAt time.Time `json:"created_at"`
// }
