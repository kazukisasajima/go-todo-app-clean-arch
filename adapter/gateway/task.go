package gateway

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"go-todo-app-clean-arch/entity"
)

type TaskRepository interface {
	Create(task *entity.Task) (*entity.Task, error)
	Get(userId int, taskId int) (*entity.Task, error)
	GetAllTasks(userId int) ([]*entity.Task, error)
	Save(task *entity.Task, userId int, taskId int) (*entity.Task, error)
	Delete(userId int, taskId int) error
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

func (t *taskRepository) Get(userId int, taskId int) (*entity.Task, error) {
	task := entity.Task{}
	if err := t.db.
		Where("user_id = ? AND id = ?", userId, taskId).
		First(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (t *taskRepository) GetAllTasks(userId int) ([]*entity.Task, error) {
	var tasks []*entity.Task
	if err := t.db.Where("user_id = ?", userId).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskRepository) Save(task *entity.Task, userId int, taskId int) (*entity.Task, error) {
	selectedTask, err := t.Get(userId, taskId)
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

func (t *taskRepository) Delete(taskId int, userId int) error {
	task := entity.Task{ID: taskId, UserID: userId}
	if err := t.db.Where("id = ? AND user_id=?", taskId, userId).Delete(&task).Error; err != nil {
		return err
	}
	return nil
}
