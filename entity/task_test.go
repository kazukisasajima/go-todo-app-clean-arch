package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-todo-app-clean-arch/entity"
)

func TestTask(t *testing.T) {
	task := entity.Task{
		ID:        1,
		Title:     "Test Task",
		UserID:    1,
	}

	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, 1, task.UserID)
}
