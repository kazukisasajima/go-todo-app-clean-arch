package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-todo-app-clean-arch/entity"
)

func TestUser(t *testing.T) {
	user := entity.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "password123",
	}

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "password123", user.Password)
}
