package usecase

import (
	"go-todo-app-clean-arch/adapter/gateway"
	"go-todo-app-clean-arch/entity"
)

type TaskUseCase interface {
	Create(task *entity.Task) (*entity.Task, error)
	Get(userId int, taskId int) (*entity.Task, error)
	GetAllTasks(userId int) ([]*entity.Task, error)
	Save(task *entity.Task, userId int, taskId int) (*entity.Task, error)
	Delete(taskId int, userId int) error
}

type taskUseCase struct {
	taskRepository gateway.TaskRepository
}

func NewTaskUseCase(taskRepository gateway.TaskRepository) *taskUseCase {
	return &taskUseCase{
		taskRepository: taskRepository,
	}
}

func (t *taskUseCase) Create(task *entity.Task) (*entity.Task, error) {
	return t.taskRepository.Create(task)
}

func (t *taskUseCase) Get(userId int, task_id int) (*entity.Task, error) {
	return t.taskRepository.Get(userId, task_id)
}

func (t *taskUseCase) GetAllTasks(userId int) ([]*entity.Task, error) {
	return t.taskRepository.GetAllTasks(userId)
}

func (t *taskUseCase) Save(task *entity.Task, userId int, taskId int) (*entity.Task, error) {
	return t.taskRepository.Save(task, userId, taskId)
}

func (t *taskUseCase) Delete(userId int, taskId int) error {
	return t.taskRepository.Delete(userId, taskId)
}
