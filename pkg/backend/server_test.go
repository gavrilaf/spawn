package backend

import (
	"context"
	"testing"

	"github.com/gavrilaf/spawn/pkg/senv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartServer(t *testing.T) {
	srv := CreateServer(senv.GetEnvironment())
	require.NotNil(t, srv)

	assert.Equal(t, StateCreated, srv.GetServerState())

	_, err := srv.Ping(context.Background(), nil)
	assert.NotNil(t, err)

	srv.StartServer()

	assert.Equal(t, StateOk, srv.GetServerState())

	_, err = srv.Ping(context.Background(), nil)
	assert.Nil(t, err)
}
