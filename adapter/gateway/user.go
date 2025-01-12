package gateway

import (
	"gorm.io/gorm"
	"go-todo-app-clean-arch/entity"
)

type UserRepository interface {
	Signup(user *entity.User) (*entity.User, error)
	GetCurrentUser(userId int) (*entity.User, error)
	DeleteUser(userId int) error
	FindByEmail(email string) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (u *userRepository) Signup(user *entity.User) (*entity.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) GetCurrentUser(userId int) (*entity.User, error) {
	user := entity.User{}
	if err := u.db.First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) DeleteUser(userId int) error {
	if err := u.db.Delete(&entity.User{}, userId).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepository) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	if err := u.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
