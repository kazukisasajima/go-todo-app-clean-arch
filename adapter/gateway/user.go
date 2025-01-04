package gateway

import (
	"gorm.io/gorm"
	"go-todo-app-clean-arch/entity"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Get(ID int) (*entity.User, error)
	Delete(ID int) error
	GetTasksForUser(userID int) ([]*entity.Task, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (u *userRepository) Create(user *entity.User) (*entity.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) Get(ID int) (*entity.User, error) {
	user := entity.User{}
	if err := u.db.First(&user, ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) Delete(ID int) error {
	if err := u.db.Delete(&entity.User{}, ID).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepository) GetTasksForUser(userID int) ([]*entity.Task, error) {
	var tasks []*entity.Task
	if err := u.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
