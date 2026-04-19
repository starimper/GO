package service

import (
	"errors"
	"practice-8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)


func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	expected := &repository.User{ID: 1, Name: "Ruslan"}
	mockRepo.EXPECT().GetUserByID(1).Return(expected, nil)

	result, err := svc.GetUserByID(1)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}


func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "Ruslan"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := svc.CreateUser(user)
	require.NoError(t, err)
}


func TestRegisterUser_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	existing := &repository.User{ID: 1, Name: "Old", Email: "old@kbtu.kz"}
	mockRepo.EXPECT().GetByEmail("old@kbtu.kz").Return(existing, nil)

	err := svc.RegisterUser(&repository.User{Name: "New"}, "old@kbtu.kz")
	require.Error(t, err)
	assert.ErrorContains(t, err, "already exists")
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{ID: 3, Name: "Ruslan", Email: "ruslan@kbtu.kz"}
	mockRepo.EXPECT().GetByEmail("ruslan@kbtu.kz").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(nil)

	err := svc.RegisterUser(newUser, "ruslan@kbtu.kz")
	require.NoError(t, err)
}

func TestRegisterUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetByEmail("x@kbtu.kz").Return(nil, errors.New("db connection lost"))

	err := svc.RegisterUser(&repository.User{Name: "X"}, "x@kbtu.kz")
	require.Error(t, err)
	assert.ErrorContains(t, err, "error getting user")
}


func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := NewUserService(repository.NewMockUserRepository(ctrl))

	err := svc.UpdateUserName(5, "")
	require.Error(t, err)
	assert.ErrorContains(t, err, "name cannot be empty")
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(99).Return(nil, errors.New("not found"))

	err := svc.UpdateUserName(99, "NewName")
	require.Error(t, err)
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 5, Name: "OldName"}
	mockRepo.EXPECT().GetUserByID(5).Return(user, nil)

	// DoAndReturn lets us assert the name was actually changed BEFORE UpdateUser is called.
	mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(
		func(u *repository.User) error {
			assert.Equal(t, "NewName", u.Name, "name must be mutated before UpdateUser is called")
			return nil
		},
	)

	err := svc.UpdateUserName(5, "NewName")
	require.NoError(t, err)
}

func TestUpdateUserName_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 5, Name: "OldName"}
	mockRepo.EXPECT().GetUserByID(5).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("update failed"))

	err := svc.UpdateUserName(5, "NewName")
	require.Error(t, err)
	assert.ErrorContains(t, err, "update failed")
}


func TestDeleteUser_AdminBlocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := NewUserService(repository.NewMockUserRepository(ctrl))

	err := svc.DeleteUser(1)
	require.Error(t, err)
	assert.ErrorContains(t, err, "not allowed to delete admin")
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(42).Return(nil)

	err := svc.DeleteUser(42)
	require.NoError(t, err)
}

func TestDeleteUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(42).Return(errors.New("db is down"))

	err := svc.DeleteUser(42)
	require.Error(t, err)
	assert.ErrorContains(t, err, "db is down")
}
