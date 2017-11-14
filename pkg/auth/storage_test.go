package auth

import (
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageMock(t *testing.T) {
	storage := NewStorageMock(env.GetEnvironment("Test"))

	_, err := storage.FindClient("client_test")
	assert.Nil(t, err)

	_, err = storage.FindUser("id1@spawn.com")
	assert.Nil(t, err)

	_, err = storage.FindUser("id2@spawn.com")
	assert.Nil(t, err)
}
