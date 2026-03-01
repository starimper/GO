package repository

import (
	"practice_4/internal/repository/postgres"
	"practice_4/internal/repository/postgres/users"
	"practice_4/pkg/modules"
)

type UserRepository interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user *modules.User) (int, error)
	UpdateUser(user *modules.User) error
	DeleteUser(id int) (int64, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}