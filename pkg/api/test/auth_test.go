package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	//"github.com/stretchr/testify/require"
	//"github.com/fatih/structs"
	"github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"
	"testing"
	//"github.com/gin-gonic/gin"

	//"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/cryptx"
)

func Test_SignUp(t *testing.T) {
	teng := createTestEngine(t)

	deviceID := "device-111"
	username := uuid.NewV4().String()
	client := teng.getClient(t)
	sign := cryptx.GenerateSignature(client.ID+deviceID+username, client.Secret)

	correct := map[string]string{
		"client_id":   client.ID,
		"device_id":   deviceID,
		"device_name": "Test device",
		"username":    username,
		"password":    "password",
		"signature":   string(sign)}

	invalid_param := map[string]string{
		"client_id": client.ID,
		"device_id": deviceID}

	invalid_sign := map[string]string{
		"client_id":   client.ID,
		"device_id":   deviceID,
		"device_name": "Test device",
		"username":    username,
		"password":    "password",
		"signature":   "11111"}

	already_registered := correct

	tests := []struct {
		body          map[string]string
		expected_code int
	}{
		{correct, 200},
		{invalid_param, 400},
		{invalid_sign, 500},
		{already_registered, 500},
	}

	for _, tt := range tests {
		jbody, _ := json.Marshal(tt.body)
		req, _ := http.NewRequest("PUT", "/auth/register", bytes.NewReader(jbody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		teng.engine.ServeHTTP(w, req)

		assert.Equal(t, tt.expected_code, w.Code)
	}
}

func Test_SignIn(t *testing.T) {
	teng := createTestEngine(t)

	correct := teng.registerUser(t)
	correct["auth_type"] = "password"

	invalid_param := mapDeepCopy(correct)
	delete(invalid_param, "device_id")

	invalid_sign := mapDeepCopy(correct)
	invalid_param["signature"] = "11111"

	tests := []struct {
		name          string
		body          map[string]string
		expected_code int
	}{
		{"correct", correct, 200},
		{"invalid param", invalid_param, 400},
		{"invalid sign", invalid_sign, 500},
	}

	for _, tt := range tests {
		jbody, _ := json.Marshal(tt.body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(jbody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		teng.engine.ServeHTTP(w, req)

		assert.Equal(t, tt.expected_code, w.Code, tt.name)
	}

}
