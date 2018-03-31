package auth

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/senv"
)

const (
	tClientID = "client-test-01"
	tDeviveID = "device1"
	tPsw      = "password"
)

func getApi(t *testing.T) ApiImpl {
	bridge := api.CreateBridge(senv.GetEnvironment())
	require.NotNil(t, bridge)
	return CreateApi(bridge)
}

func getClient(t *testing.T, api ApiImpl) *db.Client {
	p, err := api.getClient(tClientID)
	require.Nil(t, err)
	return p
}

/////////////////////////////////////////////////////////////////////////////////////////////////

func Test_SignUp(t *testing.T) {
	api := getApi(t)
	client := getClient(t, api)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	dto := RegisterDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID,
		Username:  username,
		Password:  tPsw,
		Signature: sign}

	token, err := api.handleSignUp(dto, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	assert.NotEmpty(t, token.AuthToken)
	assert.NotEmpty(t, token.RefreshToken)

	assert.Equal(t, false, token.Permissions.IsLocked)
	assert.Equal(t, false, token.Permissions.IsEmailConfirmed)
	assert.Equal(t, false, token.Permissions.Is2FARequired)
	assert.Equal(t, true, token.Permissions.IsDeviceConfirmed)

	// Already registered
	_, err = api.handleSignUp(dto, LoginContext{})
	require.Equal(t, types.ErrUserAlreadyExist, err)

	// Invalid signature
	dto.Signature += "111"
	_, err = api.handleSignUp(dto, LoginContext{})
	require.Equal(t, types.ErrInvalidSignature, err)
}

func Test_SignIn(t *testing.T) {
	api := getApi(t)
	client := getClient(t, api)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	regDto := RegisterDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID,
		Username:  username,
		Password:  tPsw,
		Signature: sign}

	_, err := api.handleSignUp(regDto, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	sign2 := cryptx.GenerateSignature(client.ID+tDeviveID+"111"+username, client.Secret)
	logingDto := LoginDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID + "111",
		AuthType:  types.AuthTypeSimple,
		Username:  username,
		Password:  tPsw,
		Signature: sign2}

	token, err := api.handleSignIn(logingDto, LoginContext{IP: "127.0.0.1", Region: "SF"})
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
	_, err = api.handleSignIn(logingDto, LoginContext{})
	assert.Equal(t, types.ErrUserUnknown, err)

	logingDto.Signature += "111"
	_, err = api.handleSignIn(logingDto, LoginContext{})
	require.Equal(t, types.ErrInvalidSignature, err)
}

func Test_RefreshToken(t *testing.T) {
	api := getApi(t)
	client := getClient(t, api)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	dto := RegisterDTO{
		ClientID:  client.ID,
		DeviceID:  tDeviveID,
		Username:  username,
		Password:  tPsw,
		Signature: sign}

	token, err := api.handleSignUp(dto, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	refreshDto := RefreshDTO{
		AuthToken:    token.AuthToken,
		RefreshToken: token.RefreshToken,
	}

	auth, err := api.handleRefresh(refreshDto)
	assert.Nil(t, err)
	require.NotNil(t, auth)

	assert.NotEmpty(t, auth.AuthToken)
	assert.NotEqual(t, token.AuthToken, auth.AuthToken)
	assert.Empty(t, auth.RefreshToken)
	assert.Equal(t, token.Permissions, auth.Permissions)

	assert.Equal(t, float64(1), auth.Expire.Sub(time.Now()).Round(time.Hour).Hours()) // One hour token
}
