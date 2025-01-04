package gateway

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"go-todo-app-clean-arch/entity"
)

type TaskRepository interface {
	Create(task *entity.Task) (*entity.Task, error)
	Get(ID int) (*entity.Task, error)
	Save(*entity.Task) (*entity.Task, error)
	Delete(ID int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (t *taskRepository) Create(task *entity.Task) (*entity.Task, error) {
	if err := t.db.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskRepository) Get(ID int) (*entity.Task, error) {
	task := entity.Task{}
	if err := t.db.First(&task, ID).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (t *taskRepository) Save(task *entity.Task) (*entity.Task, error) {
	selectedTask, err := t.Get(task.ID)
	if err != nil {
		return nil, err
	}

	if err := copier.CopyWithOption(selectedTask, task, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}
	if err := t.db.Save(selectedTask).Error; err != nil {
		return nil, err
	}

	return selectedTask, nil
}

func (t *taskRepository) Delete(ID int) error {
	task := entity.Task{ID: ID}
	if err := t.db.Where("id = ?", &task.ID).Delete(&task).Error; err != nil {
		return err
	}
	return nil
}
