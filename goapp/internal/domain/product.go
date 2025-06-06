package domain

type Product struct {
	ID          int     `json:"id"`
	Code        string  `json:"codigo"`
	Description string  `json:"descripcion"`
	Price       float64 `json:"precio"`
	Stock       int     `json:"stock"`
}

type ProductRepository interface {
	GetAll() ([]Product, error)
	Create(Product) error
}
