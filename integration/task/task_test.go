package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"go-todo-app-clean-arch/adapter/controller/gin/presenter"
	"go-todo-app-clean-arch/pkg"
)

type TaskTestSuite struct {
	suite.Suite
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}

func (suite *TaskTestSuite) TestTaskCreateGetDelete() {
	// Create
	baseEndpoint := pkg.GetEndpoint("/api/v1")
	apiClient, _ := presenter.NewClientWithResponses(baseEndpoint)
	createResponse, err := apiClient.CreateTaskWithResponse(context.Background(), presenter.CreateTaskJSONRequestBody{
		Title: "test title",
	})
	suite.Assert().Nil(err)
	suite.Equal(http.StatusCreated, createResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().NotNil(createResponse.JSON201.Id)
	suite.Assert().Equal("test title", createResponse.JSON201.Title)

	// Get
	getResponse, err := apiClient.GetTaskByIdWithResponse(context.Background(), createResponse.JSON201.Id)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, getResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().Equal(createResponse.JSON201.Id, getResponse.JSON200.Id)
	suite.Assert().Equal("test title", getResponse.JSON200.Title)

	// Update
	title := "test title updated"
	updateResponse, err := apiClient.UpdateTaskByIdWithResponse(context.Background(), getResponse.JSON200.Id, presenter.UpdateTaskByIdJSONRequestBody{
		Title: &title,
	})
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, updateResponse.StatusCode())
	suite.Assert().Nil(err)
	suite.Assert().Equal("test title updated", updateResponse.JSON200.Title)
	
	// Delete
	deleteResponse, err := apiClient.DeleteTaskByIdWithResponse(context.Background(), updateResponse.JSON200.Id)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusNoContent, deleteResponse.StatusCode())
}