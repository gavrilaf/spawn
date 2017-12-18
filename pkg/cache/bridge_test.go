package cache

import (
	//"fmt"

	"testing"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/stretchr/testify/require"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/require"
)

func TestBridge_Connect(t *testing.T) {

	cache, err := Connect(env.GetEnvironment("Test"))
	defer cache.Close()

	require.Nil(t, err)
	require.NotNil(t, cache)

	br := cache.(*Bridge)

	conn := br.get()
	defer conn.Close()

	err = conn.Err()

	require.Nil(t, err)
}
