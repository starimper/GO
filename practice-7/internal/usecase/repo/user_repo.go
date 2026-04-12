package repo

import (
	"fmt"
	"practice-7/internal/entity"
	"practice-7/pkg/postgres"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{PG: pg}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	if err := u.PG.Conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) LoginUser(dto *entity.LoginUserDTO) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("username = ?", dto.Username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("username not found: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) GetUserByID(userID string) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) PromoteUser(userID string) error {
	result := u.PG.Conn.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("role", "admin")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", userID)
	}
	return nil
}
