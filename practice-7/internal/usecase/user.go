package usecase

import (
	"fmt"
	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	created, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}
	sessionID := uuid.New().String()
	return created, sessionID, nil
}

func (u *UserUseCase) LoginUser(dto *entity.LoginUserDTO) (string, error) {
	userFromRepo, err := u.repo.LoginUser(dto)
	if err != nil {
		return "", fmt.Errorf("login – repo: %w", err)
	}

	if !utils.CheckPassword(userFromRepo.Password, dto.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("login – generate JWT: %w", err)
	}

	return token, nil
}


func (u *UserUseCase) GetMe(userID string) (*entity.User, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("get me: %w", err)
	}
	return user, nil
}


func (u *UserUseCase) PromoteUser(userID string) error {
	if err := u.repo.PromoteUser(userID); err != nil {
		return fmt.Errorf("promote user: %w", err)
	}
	return nil
}
