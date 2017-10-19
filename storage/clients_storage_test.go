package storage

import (
	"fmt"
	"testing"
)

func TestClientStorage(t *testing.T) {
	strg := NewClientsStorage()

	key, err := strg.FindClient("test")
	if err != nil {
		t.Fatal("Can't fing client")
	}

	fmt.Printf("Found client: %v\n", key)
}
