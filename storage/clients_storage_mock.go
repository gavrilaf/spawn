package storage

import (
	"fmt"
)

type ClientsStorageMock struct{}

var DefaultClients = map[string][]byte{
	"client_test": []byte("9adfb490d6b7ea759f56875c89b5db6e7850b1638e193694481294d01f098575"),
	"client_ios":  []byte("d81d3e25c0f83c6ea0efcde45cad98b3501ec3f21ae01605499e95b77a4a3366"),
}

func (c ClientsStorageMock) FindClientByID(id string) (*Client, error) {
	secret, ok := DefaultClients[id]
	if !ok {
		return nil, fmt.Errorf("Can't find client by id %v", id)
	}
	return &Client{ID: id, Secret: secret}, nil
}
