package domain

type Product struct {
	ID          int     `json:"id" gorm:"column:idproducto;primaryKey"` // Assuming idproducto is PK for 'productos' table
	Code        string  `json:"codigo" gorm:"column:codigo"`
	Description string  `json:"descripcion" gorm:"column:descripcion"`
	Price       float64 `json:"precio" gorm:"column:precio"`
	Stock       int     `json:"stock" gorm:"column:existencia"` // PHP side often uses 'existencia' for stock
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	Create(Product) error
}
