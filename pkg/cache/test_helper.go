package cache

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/senv"
	"github.com/stretchr/testify/require"
)

func getTestCache(t *testing.T) Cache {
	cache := Connect(senv.GetEnvironment())
	require.NotNil(t, cache)
	return cache
}
