package test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/config"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/senv"
)

const (
	tClientID = "client-test-01"
)

type testEngine struct {
	api    *api.Bridge
	engine *gin.Engine
}

func createTestEngine(t *testing.T) testEngine {
	env := senv.GetEnvironment()
	bridge := api.CreateBridge(env)
	require.NotNil(t, bridge)

	engine := gin.Default()
	config.ConfigureEngine(engine, bridge)

	return testEngine{api: bridge, engine: engine}
}

func (self testEngine) getClient(t *testing.T) *db.Client {
	p, err := self.api.ReadModel.FindClient(tClientID)
	require.Nil(t, err)
	return p
}
