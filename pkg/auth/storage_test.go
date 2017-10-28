package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageWithMock(t *testing.T) {
	storage := StorageFacade{Clients: NewClientsStorageMock(), Users: NewUsersStorageMock(), Sessions: NewMemorySessionsStorage()}

	client, err := storage.FindClientByID("client_test")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	user, err := storage.FindUserByUsername("id1@i.com")
	assert.Nil(t, err)
	assert.NotNil(t, user)

}
