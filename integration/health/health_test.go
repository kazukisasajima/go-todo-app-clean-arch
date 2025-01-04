package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-todo-app-clean-arch/pkg"
)

func TestPing(t *testing.T) {
	endpoint := pkg.GetEndpoint("/health")
	res, err := http.Get(endpoint)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}
