package cache

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/require"
)

func getTestCache(t *testing.T) Cache {
	cache, err := Connect(env.GetEnvironment("Test"))
	require.Nil(t, err)
	require.NotNil(t, cache)
	return cache
}
