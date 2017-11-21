package backend

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/env"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartServer(t *testing.T) {
	srv := CreateServer(env.GetEnvironment("Test"))
	require.NotNil(t, srv)

	assert.Equal(t, StateCreated, srv.state)

	srv.StartServer()

	assert.Equal(t, StateOk, srv.state)
}
