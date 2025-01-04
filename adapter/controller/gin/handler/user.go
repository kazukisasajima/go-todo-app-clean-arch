package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-todo-app-clean-arch/adapter/controller/gin/presenter"
	"go-todo-app-clean-arch/entity"
	"go-todo-app-clean-arch/pkg/logger"
	"go-todo-app-clean-arch/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func userToResponse(user *entity.User) *presenter.UserResponse {
	return &presenter.UserResponse{
		Id:    user.ID,
		Email: user.Email,
	}
}

func (u *UserHandler) CreateUser(c *gin.Context) {
	var requestBody presenter.CreateUserJSONRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	user := &entity.User{
		Email:    string(requestBody.Email),
		Password: requestBody.Password,
	}

	createdUser, err := u.userUseCase.Create(user)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userToResponse(createdUser))
}

func (u *UserHandler) GetUserById(c *gin.Context, ID int) {
	user, err := u.userUseCase.Get(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, userToResponse(user))
}

func (u *UserHandler) DeleteUserById(c *gin.Context, ID int) {
	err := u.userUseCase.Delete(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (u *UserHandler) GetTasksForUser(c *gin.Context, ID int) {
	tasks, err := u.userUseCase.GetTasksForUser(ID)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
		return
	}

	var response []presenter.TaskResponse
	for _, task := range tasks {
		response = append(response, presenter.TaskResponse{
			Id:     task.ID,
			Title:  task.Title,
			UserId: task.UserID,
		})
	}

	c.JSON(http.StatusOK, response)
}
