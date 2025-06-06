package mysql

import (
	"database/sql"
	"opticav2/internal/domain"
)

type ProductRepository struct {
	DB *sql.DB
}

func (r ProductRepository) GetAll() ([]domain.Product, error) {
	rows, err := r.DB.Query("SELECT codproducto, codigo, descripcion, precio, existencia FROM producto")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Code, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r ProductRepository) Create(p domain.Product) error {
	_, err := r.DB.Exec("INSERT INTO producto(codigo, descripcion, precio, existencia) VALUES(?,?,?,?)",
		p.Code, p.Description, p.Price, p.Stock)
	return err
}
