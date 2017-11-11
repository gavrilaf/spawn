package auth

import (
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageMock(t *testing.T) {
	storage := NewStorageMock(env.GetEnvironment("Test"))

	client, err := storage.FindClient("client_test")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	user, err := storage.FindUser("id1@spawn.com")
	assert.Nil(t, err)
	assert.NotNil(t, user)
}
