package storage

import (
	"fmt"
	types "github.com/gavrilaf/go-auth"
)

type ClientsStorageImpl struct{}

func NewClientsStorage() types.ClientsStorage {
	return &ClientsStorageImpl{}
}

func (c *ClientsStorageImpl) FindClient(clientId string) (*types.ClientKey, error) {
	switch clientId {
	case "test":
		return &types.ClientKey{"test", "1234567"}, nil
	}
	return nil, fmt.Errorf("Unknown client: %v", clientId)
}
