package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-todo-app-clean-arch/adapter/controller/echo/presenter"
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

func (u *UserHandler) CreateUser(c echo.Context) error {
	var requestBody presenter.CreateUserJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	user := &entity.User{
		Email:    string(requestBody.Email),
		Password: requestBody.Password,
	}

	createdUser, err := u.userUseCase.Create(user)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, userToResponse(createdUser))
}

func (u *UserHandler) GetUserById(c echo.Context, ID int) error {
	user, err := u.userUseCase.Get(ID)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, userToResponse(user))
}

func (u *UserHandler) DeleteUserById(c echo.Context, ID int) error {
	if err := u.userUseCase.Delete(ID); err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (u *UserHandler) GetTasksForUser(c echo.Context, ID int) error {
	tasks, err := u.userUseCase.GetTasksForUser(ID)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	var response []presenter.TaskResponse
	for _, task := range tasks {
		response = append(response, presenter.TaskResponse{
			Id:     task.ID,
			Title:  task.Title,
			UserId: task.UserID,
		})
	}

	return c.JSON(http.StatusOK, response)
}
