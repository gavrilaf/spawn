package storage

import (
	"fmt"
)

type ClientsStorageMock struct{}

var DefaultClients = map[string]string{
	"client_test": "9adfb490d6b7ea759f56875c89b5db6e7850b1638e193694481294d01f098575",
	"client_ios":  "d81d3e25c0f83c6ea0efcde45cad98b3501ec3f21ae01605499e95b77a4a3366",
}

func (c ClientsStorageMock) FindClient(clientId string) (*ClientKey, error) {
	secret, ok := DefaultClients[clientId]
	if !ok {
		return nil, fmt.Errorf("Can't find secret for clientId %v", clientId)
	}
	return &ClientKey{ClientID: clientId, Secret: secret}, nil
}
