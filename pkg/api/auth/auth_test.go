package auth

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/env"
	//"github.com/stretchr/testify/assert"
	"time"

	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOnRealEnvironment = true

const (
	tClientID = "client-test-01"
	tDeviveID = "device1"
	tPsw      = "password"
)

func getMiddleware(t *testing.T) *Middleware {
	bridge := api.CreateBridge(env.GetEnvironment("Test"))
	require.NotNil(t, bridge)
	return CreateMiddleware(bridge)
}

func getClient(t *testing.T, mw *Middleware) *db.Client {
	p, err := mw.getClient(tClientID)
	require.Nil(t, err)
	return p
}

/////////////////////////////////////////////////////////////////////////////////////////////////

func TestAuth_Register(t *testing.T) {
	mw := CreateMiddleware(api.CreateBridge(env.GetEnvironment("Test")))
	require.NotNil(t, mw)

	client := getClient(t, mw)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	dto := RegisterDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID,
		Username:  username,
		Password:  tPsw,
		Signature: sign}

	token, err := mw.HandleRegister(dto, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	assert.NotEmpty(t, token.AuthToken)
	assert.NotEmpty(t, token.RefreshToken)

	assert.Equal(t, false, token.Permissions.IsLocked)
	assert.Equal(t, false, token.Permissions.IsEmailConfirmed)
	assert.Equal(t, false, token.Permissions.Is2FARequired)
	assert.Equal(t, true, token.Permissions.IsDeviceConfirmed)

	// Already registered
	_, err = mw.HandleRegister(dto, LoginContext{})
	require.Equal(t, api.ErrUserAlreadyExist, err)

	// Invalid signature
	dto.Signature += "111"
	_, err = mw.HandleRegister(dto, LoginContext{})
	require.Equal(t, api.ErrInvalidSignature, err)
}

func TestAuth_Login(t *testing.T) {
	mw := CreateMiddleware(api.CreateBridge(env.GetEnvironment("Test")))
	require.NotNil(t, mw)

	client := getClient(t, mw)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	regDto := RegisterDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID,
		Username:  username,
		Password:  tPsw,
		Signature: sign}

	_, err := mw.HandleRegister(regDto, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	sign2 := cryptx.GenerateSignature(client.ID+tDeviveID+"111"+username, client.Secret)
	logingDto := LoginDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID + "111",
		AuthType:  AuthTypeSimple,
		Username:  username,
		Password:  tPsw,
		Signature: sign2}

	token, err := mw.HandleLogin(logingDto, LoginContext{IP: "127.0.0.1", Region: "SF"})
	assert.Nil(t, err)
	assert.NotNil(t, token)

	assert.NotEmpty(t, token.AuthToken)
	assert.NotEmpty(t, token.RefreshToken)

	assert.Equal(t, false, token.Permissions.IsLocked)
	assert.Equal(t, false, token.Permissions.IsEmailConfirmed)
	assert.Equal(t, false, token.Permissions.Is2FARequired)
	assert.Equal(t, false, token.Permissions.IsDeviceConfirmed) // new device, need confirmation

	assert.Equal(t, float64(1), token.Expire.Sub(time.Now()).Round(time.Hour).Hours()) // One hour token

	logingDto.Username += "111"
	logingDto.Signature = cryptx.GenerateSignature(logingDto.ClientID+logingDto.DeviceID+logingDto.Username, client.Secret)
	_, err = mw.HandleLogin(logingDto, LoginContext{})
	assert.Equal(t, api.ErrUserUnknown, err)

	logingDto.Signature += "111"
	_, err = mw.HandleLogin(logingDto, LoginContext{})
	require.Equal(t, api.ErrInvalidSignature, err)
}

func TestAuth_Refresh(t *testing.T) {
	mw := CreateMiddleware(api.CreateBridge(env.GetEnvironment("Test")))
	require.NotNil(t, mw)

	client := getClient(t, mw)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	dto := RegisterDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID,
		Username:  username,
		Password:  tPsw,
		Signature: sign}

	token, err := mw.HandleRegister(dto, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	refreshDto := RefreshDTO{
		AuthToken:    token.AuthToken,
		RefreshToken: token.RefreshToken,
	}

	auth, err := mw.HandleRefresh(refreshDto)
	assert.Nil(t, err)
	require.NotNil(t, auth)

	assert.NotEmpty(t, auth.AuthToken)
	assert.NotEqual(t, token.AuthToken, auth.AuthToken)
	assert.Empty(t, auth.RefreshToken)
	assert.Equal(t, token.Permissions, auth.Permissions)

	assert.Equal(t, float64(1), auth.Expire.Sub(time.Now()).Round(time.Hour).Hours()) // One hour token
}
