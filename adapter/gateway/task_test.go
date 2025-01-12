package gateway_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"go-todo-app-clean-arch/adapter/gateway"
	"go-todo-app-clean-arch/entity"
	"go-todo-app-clean-arch/pkg/tester"
)

type TaskRepositorySuite struct {
	tester.DBSQLiteSuite
	repository gateway.TaskRepository
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}

func (suite *TaskRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewTaskRepository(suite.DB)
}

func (suite *TaskRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewTaskRepository(mockGormDB)
	return mock
}

func (suite *TaskRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewTaskRepository(suite.DB)
}

func (suite *TaskRepositorySuite) TestTaskRepositoryCRUD() {
	task := &entity.Task{
		Title:     "Test Task",
		UserID:   1,
	}
	task, err := suite.repository.Create(task)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(task.ID)
	suite.Assert().Equal("Test Task", task.Title)
	suite.Assert().Equal(1, task.UserID)

	getTask, err := suite.repository.Get(task.ID, task.UserID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Test Task", getTask.Title)
	suite.Assert().Equal(1, getTask.UserID)

	getTask.Title = "Updated Task"
	updateTask, err := suite.repository.Save(getTask)
	suite.Assert().Nil(err)
	suite.Assert().Equal("Updated Task", updateTask.Title)

	err = suite.repository.Delete(updateTask.ID)
	suite.Assert().Nil(err)
	deleteTask, err := suite.repository.Get(updateTask.ID, 1)
	suite.Assert().Nil(deleteTask)
	suite.Assert().Equal("record not found", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `tasks` (`title`,`user_id`) VALUES (?,?)")).
		WithArgs("Fail Task", 1).
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	task := &entity.Task{Title: "Fail Task", UserID: 1}
	createdTask, err := suite.repository.Create(task)
	suite.Assert().Nil(createdTask)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `tasks` WHERE `tasks`.`id` = ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()

	err := suite.repository.Delete(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE `tasks`.`id` = ? ORDER BY `tasks`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("get error"))

	task, err := suite.repository.Get(1, 1)
	suite.Assert().Nil(task)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tasks` WHERE `tasks`.`id` = ? ORDER BY `tasks`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("save error"))

	task := &entity.Task{ID: 1, Title: "Fail Save"}
	savedTask, err := suite.repository.Save(task)
	suite.Assert().Nil(savedTask)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("save error", err.Error())
}

