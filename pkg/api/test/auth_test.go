package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"

	//"github.com/stretchr/testify/require"
	//"github.com/fatih/structs"
	"github.com/satori/go.uuid"
	"testing"
	//"github.com/gin-gonic/gin"

	//"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	//"github.com/gavrilaf/spawn/pkg/api/config"
	//"github.com/gavrilaf/spawn/pkg/senv"
)

func Test_SignUp(t *testing.T) {
	engine := createTestEngine(t)

	deviceID := "device-111"
	username := uuid.NewV4().String()
	client := engine.getClient(t)
	sign := cryptx.GenerateSignature(client.ID+deviceID+username, client.Secret)

	correct_boby := map[string]string{
		"client_id":   client.ID,
		"device_id":   deviceID,
		"device_name": "Test device",
		"username":    username,
		"password":    "password",
		"signature":   string(sign)}

	tests := []struct {
		body          map[string]string
		expected_code int
	}{
		{correct_boby, 200},
	}

	for _, tt := range tests {
		jbody, _ := json.Marshal(tt.body)
		req, _ := http.NewRequest("PUT", "/auth/register", bytes.NewReader(jbody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		engine.engine.ServeHTTP(w, req)

		assert.Equal(t, tt.expected_code, w.Code)
	}
}
