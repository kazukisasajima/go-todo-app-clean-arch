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

type UserRepositorySuite struct {
	tester.DBSQLiteSuite
	repository gateway.UserRepository
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (suite *UserRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewUserRepository(suite.DB)
}

func (suite *UserRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewUserRepository(mockGormDB)
	return mock
}

func (suite *UserRepositorySuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewUserRepository(suite.DB)
}

func (suite *UserRepositorySuite) TestUserRepositoryCRUD() {
	user := &entity.User{
		Email:    "test@example.com",
		Password: "password",
	}
	createdUser, err := suite.repository.Create(user)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(createdUser.ID)
	suite.Assert().Equal("test@example.com", createdUser.Email)

	getUser, err := suite.repository.Get(createdUser.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test@example.com", getUser.Email)

	err = suite.repository.Delete(createdUser.ID)
	suite.Assert().Nil(err)
	deletedUser, err := suite.repository.Get(createdUser.ID)
	suite.Assert().Nil(deletedUser)
	suite.Assert().Equal("record not found", err.Error())
}

func (suite *UserRepositorySuite) TestUserCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`email`,`password`) VALUES (?,?)")).
		WithArgs("fail@example.com", "password").
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	user := &entity.User{Email: "fail@example.com", Password: "password"}
	createdUser, err := suite.repository.Create(user)
	suite.Assert().Nil(createdUser)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *UserRepositorySuite) TestUserDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).WithArgs(1).
		WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()

	err := suite.repository.Delete(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}

func (suite *UserRepositorySuite) TestUserGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnError(errors.New("get error"))

	user, err := suite.repository.Get(1)
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}
