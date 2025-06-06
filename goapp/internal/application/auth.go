package application

import "opticav2/internal/domain"

type AuthService struct {
	Repo domain.UserRepository
}

func (s AuthService) Login(user, pass string) (*domain.User, error) {
	return s.Repo.GetByCredentials(user, pass)
}
