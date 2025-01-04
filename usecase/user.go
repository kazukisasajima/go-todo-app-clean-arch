package usecase

import (
	"go-todo-app-clean-arch/adapter/gateway"
	"go-todo-app-clean-arch/entity"
)

type (
	UserUseCase interface {
		Create(user *entity.User) (*entity.User, error)
		Get(ID int) (*entity.User, error)
		Delete(ID int) error
		GetTasksForUser(userID int) ([]*entity.Task, error)
	}
)

type userUseCase struct {
	userRepository gateway.UserRepository
}

func NewUserUseCase(userRepository gateway.UserRepository) *userUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (u *userUseCase) Create(user *entity.User) (*entity.User, error) {
	return u.userRepository.Create(user)
}

func (u *userUseCase) Get(ID int) (*entity.User, error) {
	return u.userRepository.Get(ID)
}

func (u *userUseCase) Delete(ID int) error {
	return u.userRepository.Delete(ID)
}

func (u *userUseCase) GetTasksForUser(userID int) ([]*entity.Task, error) {
	return u.userRepository.GetTasksForUser(userID)
}

