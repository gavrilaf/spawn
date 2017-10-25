package storage

import (
	"fmt"
	"testing"
)

func TestStorageWithMock(t *testing.T) {
	storage := StorageFacade{Clients: ClientsStorageMock{}, Users: UsersStorageMock{}, Sessions: NewMemorySessionsStorage()}

	client, err := storage.FindClientByID("client_test")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("Client %v\n", client)

	user, err := storage.FindUserByUsername("id1@i.com")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("User %v\n", user)

}
