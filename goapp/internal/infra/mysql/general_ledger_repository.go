package mysql

import (
	"fmt"
	"opticav2/internal/domain"
	"time"

	"gorm.io/gorm"
)

type GeneralLedgerRepository struct {
	DB *gorm.DB
}

func NewGeneralLedgerRepository(db *gorm.DB) domain.GeneralLedgerRepository {
	return &GeneralLedgerRepository{DB: db}
}

// GetEntriesByDateRange fetches entries from 'ingresos' or 'egresos' tables.
func (r *GeneralLedgerRepository) GetEntriesByDateRange(entryType string, fromDate, toDate time.Time) ([]domain.GeneralLedgerEntry, error) {
	var entries []domain.GeneralLedgerEntry
	var tableName string
	amountColumnName := "monto" // Generic name in struct, specific in query

	switch entryType {
	case "ingreso":
		tableName = "ingresos"
		amountColumnName = "ingresos" // Actual column name in 'ingresos' table
	case "egreso":
		tableName = "egresos"
		amountColumnName = "egresos" // Actual column name in 'egresos' table
	default:
		return nil, fmt.Errorf("invalid entry type: %s. Must be 'ingreso' or 'egreso'", entryType)
	}

	// Adjust column selection based on actual table structure.
	// The domain.GeneralLedgerEntry struct has generic field names.
	// We map them to specific table columns here.
	// Assuming 'id_ingreso' or 'id_egreso' as primary keys.
	// Assuming 'id_usuario', 'id_cliente', 'id_venta' might not exist or be nullable in these tables.
	// If they exist, they should be selected. For now, we select common fields.
	err := r.DB.Table(tableName).
		Select(fmt.Sprintf("id, fecha, descripcion, %s as monto, id_usuario, id_cliente, id_venta", amountColumnName)).
		Where("fecha BETWEEN ? AND ?", fromDate, toDate).
		Find(&entries).Error

	if err != nil {
		return nil, err
	}

	// Set the Type field for each entry
	for i := range entries {
		if entryType == "ingreso" {
			entries[i].Type = "IngresoGeneral"
		} else {
			entries[i].Type = "EgresoGeneral"
		}
	}
	return entries, nil
}

// GetSumByDateRange calculates the sum of amounts from 'ingresos' or 'egresos' tables.
func (r *GeneralLedgerRepository) GetSumByDateRange(entryType string, fromDate, toDate time.Time) (float64, error) {
	var sum float64
	var tableName string
	var amountColumnName string

	switch entryType {
	case "ingreso":
		tableName = "ingresos"
		amountColumnName = "ingresos"
	case "egreso":
		tableName = "egresos"
		amountColumnName = "egresos"
	default:
		return 0, fmt.Errorf("invalid entry type: %s. Must be 'ingreso' or 'egreso'", entryType)
	}

	err := r.DB.Table(tableName).
		Where("fecha BETWEEN ? AND ?", fromDate, toDate).
		Select(fmt.Sprintf("COALESCE(SUM(%s), 0)", amountColumnName)). // COALESCE to handle NULL sum if no records
		Row().Scan(&sum)

	if err != nil {
		return 0, err
	}
	return sum, nil
}
