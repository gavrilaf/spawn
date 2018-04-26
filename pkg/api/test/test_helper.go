package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/config"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/senv"
)

func mapDeepCopy(p map[string]string) map[string]string {
	new := make(map[string]string)
	for k, v := range p {
		new[k] = v
	}
	return new
}

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

func (self testEngine) registerUser(t *testing.T) (map[string]string, auth.AuthTokenDTO) {
	deviceID := "device-111"
	username := uuid.NewV4().String()
	client := self.getClient(t)
	sign := cryptx.GenerateSignature(client.ID+deviceID+username, client.Secret)

	body := map[string]string{
		"client_id":   client.ID,
		"device_id":   deviceID,
		"device_name": "Test device",
		"username":    username,
		"password":    "password",
		"signature":   string(sign)}

	jbody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "/auth/register", bytes.NewReader(jbody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	self.engine.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	var authToken auth.AuthTokenDTO
	err := json.Unmarshal(w.Body.Bytes(), &authToken)
	require.Nil(t, err)

	return body, authToken
}
