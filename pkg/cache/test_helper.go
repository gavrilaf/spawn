package cache

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/stretchr/testify/require"
)

func getTestCache(t *testing.T) Cache {
	cache, err := Connect(senv.GetEnvironment())
	require.Nil(t, err)
	require.NotNil(t, cache)
	return cache
}
