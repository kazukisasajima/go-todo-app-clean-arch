package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/oapi-codegen/runtime/types"

	"go-todo-app-clean-arch/adapter/controller/gin/presenter"
	"go-todo-app-clean-arch/entity"
)

type MockUserUseCase struct {
	mock.Mock
}

func NewMockUserUseCase() *MockUserUseCase {
	return &MockUserUseCase{}
}

func (m *MockUserUseCase) Create(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) Get(ID int) (*entity.User, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) Delete(ID int) error {
	args := m.Called(ID)
	return args.Error(0)
}

func (m *MockUserUseCase) GetTasksForUser(userID int) ([]*entity.Task, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

type UserHandlerSuite struct {
	suite.Suite
	userHandler *UserHandler
}

func TestUserHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}

func (suite *UserHandlerSuite) TestCreateUser() {
	email := "test@example.com"
	password := "password"
	mockUseCase := NewMockUserUseCase()

	user := &entity.User{
		Email:    email,
		Password: password,
	}

	mockUseCase.On("Create", user).Return(&entity.User{
		ID:    1,
		Email: email,
	}, nil)
	suite.userHandler = NewUserHandler(mockUseCase)

	request, _ := presenter.NewCreateUserRequest("/api/v1", presenter.CreateUserJSONRequestBody{
		Email:    types.Email(email),
		Password: password,
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.userHandler.CreateUser(ginContext)

	suite.Assert().Equal(http.StatusCreated, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var userGetResponse presenter.UserResponse
	err := json.Unmarshal(bodyBytes, &userGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(email, userGetResponse.Email)
}

func (suite *UserHandlerSuite) TestGetUser() {
	mockUseCase := NewMockUserUseCase()
	mockUseCase.On("Get", 1).Return(&entity.User{
		ID:    1,
		Email: "test@example.com",
	}, nil)
	suite.userHandler = NewUserHandler(mockUseCase)

	request, _ := presenter.NewGetUserByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.userHandler.GetUserById(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var userGetResponse presenter.UserResponse
	err := json.Unmarshal(bodyBytes, &userGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("test@example.com", userGetResponse.Email)
}

func (suite *UserHandlerSuite) TestDeleteUser() {
	mockUseCase := NewMockUserUseCase()
	mockUseCase.On("Delete", 1).Return(nil)
	suite.userHandler = NewUserHandler(mockUseCase)

	request, _ := presenter.NewDeleteUserByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.userHandler.DeleteUserById(ginContext, 1)

	suite.Assert().Equal(http.StatusNoContent, w.Code)
}

func (suite *UserHandlerSuite) TestGetTasksForUser() {
	mockUseCase := NewMockUserUseCase()
	tasks := []*entity.Task{
		{ID: 1, Title: "Task 1", UserID: 1},
		{ID: 2, Title: "Task 2", UserID: 1},
	}
	mockUseCase.On("GetTasksForUser", 1).Return(tasks, nil)

	suite.userHandler = NewUserHandler(mockUseCase)

	request, _ := presenter.NewGetTasksForUserRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.userHandler.GetTasksForUser(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var taskResponses []presenter.TaskResponse
	err := json.Unmarshal(bodyBytes, &taskResponses)

	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Require().Len(taskResponses, 2)
	suite.Assert().Equal("Task 1", taskResponses[0].Title)
	suite.Assert().Equal("Task 2", taskResponses[1].Title)
}

