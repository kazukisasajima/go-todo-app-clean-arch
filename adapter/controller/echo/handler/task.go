package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-todo-app-clean-arch/adapter/controller/echo/presenter"
	"go-todo-app-clean-arch/entity"
	"go-todo-app-clean-arch/pkg/logger"
	"go-todo-app-clean-arch/usecase"
)

type TaskHandler struct {
	taskUseCase usecase.TaskUseCase
}

func NewTaskHandler(taskUseCase usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

func (t *TaskHandler) CreateTask(c echo.Context) error {
	var requestBody presenter.CreateTaskJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	task := &entity.Task{
		Title:  requestBody.Title,
		UserID: requestBody.UserId,
	}

	createdTask, err := t.taskUseCase.Create(task)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, presenter.TaskResponse{
		Id:     createdTask.ID,
		Title:  createdTask.Title,
		UserId: createdTask.UserID,
	})
}

func (t *TaskHandler) GetTaskById(c echo.Context, ID int) error {
	task, err := t.taskUseCase.Get(ID)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, presenter.TaskResponse{
		Id:     task.ID,
		Title:  task.Title,
		UserId: task.UserID,
	})
}

func (t *TaskHandler) UpdateTaskById(c echo.Context, ID int) error {
	var requestBody presenter.UpdateTaskByIdJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	task := &entity.Task{
		ID:    ID,
		Title: *requestBody.Title,
	}

	updatedTask, err := t.taskUseCase.Save(task)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, presenter.TaskResponse{
		Id:     updatedTask.ID,
		Title:  updatedTask.Title,
		UserId: updatedTask.UserID,
	})
}

func (t *TaskHandler) DeleteTaskById(c echo.Context, ID int) error {
	if err := t.taskUseCase.Delete(ID); err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
