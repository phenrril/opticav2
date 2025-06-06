package mysql

import (
	"database/sql"
	"opticav2/internal/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func (r UserRepository) GetByCredentials(user, password string) (*domain.User, error) {
	u := domain.User{}
	err := r.DB.QueryRow("SELECT idusuario, nombre, usuario FROM usuario WHERE usuario=? AND clave=MD5(?) AND estado=1", user, password).Scan(&u.ID, &u.Name, &u.User)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
