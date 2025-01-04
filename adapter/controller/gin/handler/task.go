package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-todo-app-clean-arch/adapter/controller/gin/presenter"
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

func taskToResponse(task *entity.Task) *presenter.TaskResponse {
	return &presenter.TaskResponse{
		Id:        task.ID,
		Title:     task.Title,
		UserId:    task.UserID,
	}
}

func (t *TaskHandler) CreateTask(c *gin.Context) {
	var requestBody presenter.CreateTaskJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	task := &entity.Task{
		Title: requestBody.Title,
		UserID: requestBody.UserId,
	}

	createdTask, err := t.taskUseCase.Create(task)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, taskToResponse(createdTask))
}

func (t *TaskHandler) GetTaskById(c *gin.Context, ID int) {
	task, err := t.taskUseCase.Get(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskToResponse(task))
}

func (t *TaskHandler) UpdateTaskById(c *gin.Context, ID int) {
	var requestBody presenter.UpdateTaskByIdJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	task := &entity.Task{
		ID: ID, 
		Title: *requestBody.Title,
	}
	updatedTask, err := t.taskUseCase.Save(task)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskToResponse(updatedTask))
}

func (t *TaskHandler) DeleteTaskById(c *gin.Context, ID int) {
	err := t.taskUseCase.Delete(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
