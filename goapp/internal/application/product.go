package application

import "opticav2/internal/domain"

type ProductService struct {
	Repo domain.ProductRepository
}

func (s ProductService) List() ([]domain.Product, error) {
	return s.Repo.GetAll()
}

func (s ProductService) Create(p domain.Product) error {
	return s.Repo.Create(p)
}
