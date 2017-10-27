package storage

import (
	"github.com/gavrilaf/go-auth/auth/cerr"
)

var DefaultClients = map[string]string{
	"client_test": "9adfb490d6b7ea759f56875c89b5db6e7850b1638e193694481294d01f098575",
	"client_ios":  "d81d3e25c0f83c6ea0efcde45cad98b3501ec3f21ae01605499e95b77a4a3366",
}

type ClientsStorageMock struct{}

func NewClientsStorageMock() *ClientsStorageMock {
	return &ClientsStorageMock{}
}

func (c *ClientsStorageMock) FindClientByID(id string) (*Client, error) {
	secret, ok := DefaultClients[id]
	if !ok {
		return nil, cerr.ClientUnknown
	}
	return &Client{id: id, secret: secret}, nil
}
