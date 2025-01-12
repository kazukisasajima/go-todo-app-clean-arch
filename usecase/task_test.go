package usecase

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go-todo-app-clean-arch/entity"
)

type mockTaskRepository struct {
	mock.Mock
}

func NewMockTaskRepository() *mockTaskRepository {
	return new(mockTaskRepository)
}

func (m *mockTaskRepository) Create(task *entity.Task) (*entity.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *mockTaskRepository) Get(userID int, ID int) (*entity.Task, error) {
	args := m.Called(userID, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *mockTaskRepository) GetAllTasks(userID int) ([]*entity.Task, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *mockTaskRepository) Save(task *entity.Task) (*entity.Task, error) {
	args := m.Called(task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *mockTaskRepository) Delete(ID int) error {
	args := m.Called(ID)
	return args.Error(0)
}

type TaskUseCaseSuite struct {
	suite.Suite
	taskUseCase *taskUseCase
}

func TestTaskUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseSuite))
}

func (suite *TaskUseCaseSuite) SetupTest() {
}

func (suite *TaskUseCaseSuite) TestCreate() {
	title := "Test Task"
	userID := 1
	mockTaskRepository := NewMockTaskRepository()
	suite.taskUseCase = NewTaskUseCase(mockTaskRepository)

	task := &entity.Task{
		Title: title,
		UserID: userID,
	}

	mockTaskRepository.On("Create", task).Return(&entity.Task{
		ID: 1, 
		Title: title,
		UserID: userID,
	}, nil)

	task, err := suite.taskUseCase.Create(task)
	suite.Assert().Nil(err)
	suite.Assert().Equal(title, task.Title)
	suite.Assert().Equal(userID, task.UserID)
}

func (suite *TaskUseCaseSuite) TestGet() {
	taskID := 1
	title := "Test Task"
	userID := 1 // ユーザーID
	mockTaskRepository := NewMockTaskRepository()
	suite.taskUseCase = NewTaskUseCase(mockTaskRepository)

	mockTaskRepository.On("Get", taskID).Return(&entity.Task{
		ID:     taskID,
		Title:  title,
		UserID: userID,
	}, nil)

	task, err := suite.taskUseCase.Get(userID, taskID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(taskID, task.ID)
	suite.Assert().Equal(title, task.Title)
	suite.Assert().Equal(userID, task.UserID)
}

func (suite *TaskUseCaseSuite) TestSave() {
	taskID := 1
	title := "Updated Task"
	userID := 1 // ユーザーID
	mockTaskRepository := NewMockTaskRepository()
	suite.taskUseCase = NewTaskUseCase(mockTaskRepository)

	task := &entity.Task{
		ID:     taskID,
		Title:  title,
		UserID: userID,
	}

	mockTaskRepository.On("Save", task).Return(&entity.Task{
		ID:     taskID,
		Title:  title,
		UserID: userID,
	}, nil)

	updatedTask, err := suite.taskUseCase.Save(task)
	suite.Assert().Nil(err)
	suite.Assert().Equal(taskID, updatedTask.ID)
	suite.Assert().Equal(title, updatedTask.Title)
	suite.Assert().Equal(userID, updatedTask.UserID)
}

func (suite *TaskUseCaseSuite) TestDelete() {
	taskID := 1
	mockTaskRepository := NewMockTaskRepository()
	suite.taskUseCase = NewTaskUseCase(mockTaskRepository)

	mockTaskRepository.On("Delete", taskID).Return(nil)

	err := suite.taskUseCase.Delete(taskID)
	suite.Assert().Nil(err)
}

func (suite *TaskUseCaseSuite) TestGetAllTasks() {
	userID := 1
	title := "Test Task"
	mockTaskRepository := NewMockTaskRepository()
	suite.taskUseCase = NewTaskUseCase(mockTaskRepository)

	mockTaskRepository.On("GetAllTasks", userID).Return([]*entity.Task{
		{
			ID:     1,
			Title:  title,
			UserID: userID,
		},
	}, nil)

	tasks, err := suite.taskUseCase.GetAllTasks(userID)
	suite.Assert().Nil(err)
	suite.Assert().Len(tasks, 1)
	suite.Assert().Equal(title, tasks[0].Title)
	suite.Assert().Equal(userID, tasks[0].UserID)
}


