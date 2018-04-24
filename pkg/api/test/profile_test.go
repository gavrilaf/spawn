package test

import (
	//"bytes"
	"encoding/json"
	//"math"
	"net/http"
	"net/http/httptest"
	//"time"

	//"github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"testing"
)

func Test_GetProfile(t *testing.T) {
	teng := createTestEngine(t)
	token := teng.registerUser(t).token

	req, _ := http.NewRequest("GET", "/profile/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AuthToken)
	w := httptest.NewRecorder()

	teng.engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "get profile")

	var profile mdl.PersonalInfo

	err := json.Unmarshal(w.Body.Bytes(), &profile)
	require.Nil(t, err)
}
