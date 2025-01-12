package usecase

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go-todo-app-clean-arch/entity"
)

type mockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *mockUserRepository {
	return new(mockUserRepository)
}

func (m *mockUserRepository) Create(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) Get(ID int) (*entity.User, error) {
	args := m.Called(ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepository) Delete(ID int) error {
	args := m.Called(ID)
	return args.Error(0)
}

func (m *mockUserRepository) GetTasksForUser(userID int) ([]*entity.Task, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

func (m *mockUserRepository) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

type UserUseCaseSuite struct {
	suite.Suite
	userUseCase *userUseCase
}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}

func (suite *UserUseCaseSuite) TestCreate() {
	email := "test@example.com"
	password := "password123"
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	user := &entity.User{
		Email:    email,
		Password: password,
	}

	mockUserRepository.On("Create", user).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: password,
	}, nil)

	createdUser, err := suite.userUseCase.Create(user)
	suite.Assert().Nil(err)
	suite.Assert().Equal(email, createdUser.Email)
	suite.Assert().Equal(password, createdUser.Password)
}

func (suite *UserUseCaseSuite) TestGet() {
	userID := 1
	email := "test@example.com"
	password := "password123"
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	mockUserRepository.On("Get", userID).Return(&entity.User{
		ID:       userID,
		Email:    email,
		Password: password,
	}, nil)

	user, err := suite.userUseCase.Get(userID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(userID, user.ID)
	suite.Assert().Equal(email, user.Email)
	suite.Assert().Equal(password, user.Password)
}

func (suite *UserUseCaseSuite) TestDelete() {
	userID := 1
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	mockUserRepository.On("Delete", userID).Return(nil)

	err := suite.userUseCase.Delete(userID)
	suite.Assert().Nil(err)
}

func (suite *UserUseCaseSuite) TestGetTasksForUser() {
	userID := 1
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	mockUserRepository.On("GetTasksForUser", userID).Return([]*entity.Task{
		{
			ID:    1,
			Title: "Test Task",
		},
	}, nil)

	tasks, err := suite.userUseCase.GetTasksForUser(userID)
	suite.Assert().Nil(err)
	suite.Assert().Len(tasks, 1)
	suite.Assert().Equal("Test Task", tasks[0].Title)
}

func (suite *UserUseCaseSuite) TestSignup() {
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	user := &entity.User{
		Email:    email,
		Password: password,
	}

	mockUserRepository.On("Create", mock.AnythingOfType("*entity.User")).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
	}, nil)

	createdUser, err := suite.userUseCase.Signup(user)
	suite.Assert().Nil(err)
	suite.Assert().Equal(email, createdUser.Email)
	suite.Assert().True(CheckPasswordHash(password, createdUser.Password))
}

func (suite *UserUseCaseSuite) TestLogin() {
	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	mockUserRepository := NewMockUserRepository()
	suite.userUseCase = NewUserUseCase(mockUserRepository)

	credentials := &entity.Credentials{
		Email:    email,
		Password: password,
	}

	mockUserRepository.On("FindByEmail", email).Return(&entity.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
	}, nil)

	jwt, err := suite.userUseCase.Login(credentials)
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(jwt)
}
