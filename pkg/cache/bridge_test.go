package cache

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/stretchr/testify/require"
)

func TestBridge_Connect(t *testing.T) {
	cache := Connect(senv.GetEnvironment())
	defer cache.Close()

	require.NotNil(t, cache)

	br := cache.(*Bridge)

	conn := br.get()
	defer conn.Close()

	err := conn.Err()
	require.Nil(t, err)

	err = cache.HealthCheck()
	require.Nil(t, err)

	err = cache.Close()
	require.Nil(t, err)

	err = cache.HealthCheck()
	require.NotNil(t, err)
}
