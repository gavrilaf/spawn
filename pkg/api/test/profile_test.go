package test

import (
	//"bytes"
	"encoding/json"
	//"math"
	"net/http"
	"net/http/httptest"
	//"time"

	//"github.com/satori/go.uuid"

	"github.com/gavrilaf/spawn/pkg/api/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetProfile(t *testing.T) {
	teng := createTestEngine(t)
	userDesc, token := teng.registerUser(t)

	req, _ := http.NewRequest("GET", "/profile/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AuthToken)
	w := httptest.NewRecorder()

	teng.engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, "get profile")

	var userProfile profile.UserProfile

	err := json.Unmarshal(w.Body.Bytes(), &userProfile)
	require.Nil(t, err)

	assert.NotEmpty(t, userProfile.ID)

	assert.Empty(t, userProfile.PersonalInfo.FirstName)
	assert.Empty(t, userProfile.PersonalInfo.LastName)
	assert.Empty(t, userProfile.PersonalInfo.Country)
	assert.Equal(t, "1800-01-01", userProfile.PersonalInfo.BirthDate)

	assert.Equal(t, 0, userProfile.PersonalInfo.PhoneNumber.CountryCode)
	assert.Empty(t, userProfile.PersonalInfo.PhoneNumber.Number)
	assert.False(t, userProfile.PersonalInfo.PhoneNumber.IsConfirmed)

	assert.Equal(t, userDesc["username"], userProfile.AuthInfo.Username)
	assert.False(t, userProfile.AuthInfo.IsLocked)
	assert.False(t, userProfile.AuthInfo.Is2FARequired)
	assert.False(t, userProfile.AuthInfo.Is2FARequired)
	assert.Equal(t, 0, userProfile.AuthInfo.Scope)
}
