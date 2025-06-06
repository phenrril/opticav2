package domain

type User struct {
	ID   int    `json:"id"`
	Name string `json:"nombre"`
	User string `json:"usuario"`
}

type UserRepository interface {
	GetByCredentials(user, password string) (*User, error)
}
