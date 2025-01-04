package usecase

import (
	"go-todo-app-clean-arch/adapter/gateway"
	"go-todo-app-clean-arch/entity"
)


type (
	TaskUseCase interface {
		Create(task *entity.Task) (*entity.Task, error)
		Get(ID int) (*entity.Task, error)
		Save(*entity.Task) (*entity.Task, error)
		Delete(ID int) error
	}
)

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

func (t *taskUseCase) Get(ID int) (*entity.Task, error) {
	return t.taskRepository.Get(ID)
}

func (t *taskUseCase) Save(task *entity.Task) (*entity.Task, error) {
	return t.taskRepository.Save(task)
}

func (t *taskUseCase) Delete(ID int) error {
	return t.taskRepository.Delete(ID)
}
