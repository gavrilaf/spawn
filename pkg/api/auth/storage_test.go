package auth

import (
	//"fmt"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageMock(t *testing.T) {
	storage := NewStorageMock(env.GetEnvironment("Test"))

	client, err := storage.FindClient("client-test-01")
	assert.NotNil(t, client)
	assert.Nil(t, err)

	u1, err := storage.FindUser("id1@spawn.com")
	assert.NotNil(t, u1)
	assert.Nil(t, err)

	u2, err := storage.FindUser("id2@spawn.com")
	assert.NotNil(t, u2)
	assert.Nil(t, err)

	d1, err := storage.FindDevice(u1.ID, "d1")
	assert.NotNil(t, d1)
	assert.Nil(t, err)
}
