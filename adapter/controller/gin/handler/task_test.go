package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go-todo-app-clean-arch/adapter/controller/gin/presenter"
	"go-todo-app-clean-arch/entity"
)

type MockTaskUseCase struct {
	mock.Mock
}

func NewMockTaskUseCase() *MockTaskUseCase {
	return &MockTaskUseCase{}
}

func (m *MockTaskUseCase) Create(task *entity.Task) (*entity.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskUseCase) Get(ID int) (*entity.Task, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskUseCase) Save(task *entity.Task) (*entity.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskUseCase) Delete(ID int) error {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return nil
}

type TaskHandlerSuite struct {
	suite.Suite
	taskHandler *TaskHandler
}

func TestTaskHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(TaskHandlerSuite))
}

func (suite *TaskHandlerSuite) TestCreate() {
	title := "Test Task"
	user_id := 1
	mockUseCase := NewMockTaskUseCase()

	task := &entity.Task{
		Title: title,
		UserID: user_id,
	}

	mockUseCase.On("Create", task).Return(&entity.Task{
		ID: 1, 
		Title: title,
		UserID: user_id,
	}, nil)
	suite.taskHandler = NewTaskHandler(mockUseCase)

	request, _ := presenter.NewCreateTaskRequest("/api/v1", presenter.CreateTaskJSONRequestBody{
		Title: title,
		UserId: user_id,
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.taskHandler.CreateTask(ginContext)

	suite.Assert().Equal(http.StatusCreated, w.Code)
	bodyBytes, _ := io.ReadAll(w.Body)
	var taskGetResponse presenter.TaskResponse
	err := json.Unmarshal(bodyBytes, &taskGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusCreated, w.Code)
	suite.Assert().Equal(title, taskGetResponse.Title)
}

func (suite *TaskHandlerSuite) TestCreateRequestBodyFailure() {
    mockUseCase := NewMockTaskUseCase()
    suite.taskHandler = NewTaskHandler(mockUseCase)

    w := httptest.NewRecorder()
    ginContext, _ := gin.CreateTestContext(w)

    req, err := http.NewRequest("POST", "/api/v1/task", nil)
    suite.Require().Nil(err)
    req.Header.Add("Content-Type", "application/json")
    ginContext.Request = req

    suite.taskHandler.CreateTask(ginContext)

    suite.Assert().Equal(http.StatusBadRequest, w.Code)
    suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *TaskHandlerSuite) TestCreateFailure() {
	mockUseCase := NewMockTaskUseCase()
    suite.taskHandler = NewTaskHandler(mockUseCase)

	task := &entity.Task{
		Title: "Test Task",
		UserID: 1,
	}

	mockUseCase.On("Create", task).Return(nil, errors.New("invalid"))
	suite.taskHandler = NewTaskHandler(mockUseCase)

	request, _ := presenter.NewCreateTaskRequest("/api/v1", presenter.CreateTaskJSONRequestBody{
		Title: "Test Task",
		UserId: 1,
	})
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.taskHandler.CreateTask(ginContext)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message":"invalid"}`, w.Body.String())
}

func (suite *TaskHandlerSuite) TestGet() {
	mockUseCase := NewMockTaskUseCase()
	mockUseCase.On("Get", 1).Return(&entity.Task{
		ID:          1,
		Title:       "Test Task",
	}, nil)
	suite.taskHandler = NewTaskHandler(mockUseCase)

	request, _ := presenter.NewGetTaskByIdRequest("/api/v1", 1)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.taskHandler.GetTaskById(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var taskGetResponse presenter.TaskResponse
	err := json.Unmarshal(bodyBytes, &taskGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("Test Task", taskGetResponse.Title)
}

func (suite *TaskHandlerSuite) TestGetNoTaskFailure() {
	mockUseCase := NewMockTaskUseCase()
	mockUseCase.On("Get", 1111).Return(nil,
		errors.New("invalid"))
	suite.taskHandler = NewTaskHandler(mockUseCase)

	request, _ := presenter.NewGetTaskByIdRequest("/api/v1", 1111)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.taskHandler.GetTaskById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}

func (suite *TaskHandlerSuite) TestUpdate() {
	mockUseCase := NewMockTaskUseCase()
	title := "updated"

	task := &entity.Task{
		ID:       1,
		Title:    title,
	}

	mockUseCase.On("Save", task).Return(&entity.Task{
		ID:          1,
		Title:       title,
	}, nil)

	suite.taskHandler = NewTaskHandler(mockUseCase)

	request, _ := presenter.NewUpdateTaskByIdRequest("/api/v1", 1,
		presenter.UpdateTaskByIdJSONRequestBody{
			Title:    &title,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.taskHandler.UpdateTaskById(ginContext, 1)

	bodyBytes, _ := io.ReadAll(w.Body)
	var taskGetResponse presenter.TaskResponse
	err := json.Unmarshal(bodyBytes, &taskGetResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("updated", taskGetResponse.Title)
}

func (suite *TaskHandlerSuite) TestUpdateRequestBodyFailure() {
	mockUseCase := NewMockTaskUseCase()
    suite.taskHandler = NewTaskHandler(mockUseCase)

	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("PATCH", "/api/v1/task", nil)
	req.Header.Add("Content-Type", "application/json")
	ginContext.Request = req

	suite.taskHandler.CreateTask(ginContext)

	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid request"}`, w.Body.String())
}

func (suite *TaskHandlerSuite) TestUpdateTaskFailure() {
	mockUseCase := NewMockTaskUseCase()
	title := "updated"
	task := &entity.Task{
		ID:       1111,
		Title:    title,
	}

	mockUseCase.On("Save", task).Return(nil, errors.New("invalid"))
	suite.taskHandler = NewTaskHandler(mockUseCase)

	request, _ := presenter.NewUpdateTaskByIdRequest("/api/v1", 1111,
		presenter.UpdateTaskByIdJSONRequestBody{
			Title:    &title,
		},
	)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.taskHandler.UpdateTaskById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}

func (suite *TaskHandlerSuite) TestDelete() {
    mockUseCase := NewMockTaskUseCase()
    mockUseCase.On("Delete", 1).Return(nil, nil)
    suite.taskHandler = NewTaskHandler(mockUseCase)

    request, err := presenter.NewDeleteTaskByIdRequest("/api/v1", 1)
    suite.Require().Nil(err)
    w := httptest.NewRecorder()
    ginContext, _ := gin.CreateTestContext(w)
    ginContext.Request = request

    suite.taskHandler.DeleteTaskById(ginContext, 1)

    suite.Assert().Equal(http.StatusNoContent, w.Code)
}

func (suite *TaskHandlerSuite) TestDeleteTaskFailure() {
	mockUseCase := NewMockTaskUseCase()
	mockUseCase.On("Delete", 1111).Return(nil, errors.New("invalid"))
	suite.taskHandler = NewTaskHandler(mockUseCase)

	request, _ := presenter.NewDeleteTaskByIdRequest("/api/v1", 1111)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.taskHandler.DeleteTaskById(ginContext, 1111)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().JSONEq(`{"message": "invalid"}`, w.Body.String())
}