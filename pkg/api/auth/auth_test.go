package auth

import (
	"testing"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/env"
	//"github.com/stretchr/testify/assert"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/model"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testOnRealEnvironment = true

const (
	tClientID = "client_test"
	tDeviveID = "device1"
	tPsw      = "password"
)

var mockMiddleware = CreateMockMiddleware()

func GetMockMiddleware() *Middleware {
	return mockMiddleware
}

func GetMiddleware(t *testing.T) *Middleware {
	bridge := api.CreateBridge(env.GetEnvironment("Test"))
	require.NotNil(t, bridge)
	return CreateMiddleware(bridge)
}

func GetClient(t *testing.T, mw *Middleware) *db.Client {
	p, err := storageMock.FindClient(tClientID)
	require.Nil(t, err)
	return p
}

func GetRegistrationDTO(t *testing.T, mw *Middleware) RegisterDTO {
	client := GetClient(t, mw)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	return RegisterDTO{ClientID: client.ID, DeviceID: tDeviveID, Username: username, Password: tPsw, Signature: sign}
}

func GetLoginDTO(t *testing.T, reg RegisterDTO, mw *Middleware) LoginDTO {
	client := GetClient(t, mw)
	sign := cryptx.GenerateSignature(client.ID+reg.DeviceID+reg.Username, client.Secret)

	return LoginDTO{ClientID: client.ID,
		DeviceID:  reg.DeviceID,
		AuthType:  AuthTypeSimple,
		Username:  reg.Username,
		Password:  reg.Password,
		Signature: sign}
}

func DoTestRegistration(t *testing.T, mw *Middleware) {
	p := GetRegistrationDTO(t, mw)

	_, err := mw.HandleRegister(p, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	// Already registered
	_, err = mw.HandleRegister(p, LoginContext{})
	require.Equal(t, errUserAlreadyExist, err)

	// Invalid signature
	p.Signature += "111"
	_, err = mw.HandleRegister(p, LoginContext{})
	require.Equal(t, errInvalidSignature, err)
}

func DoTestLogin(t *testing.T, mw *Middleware) {
	reg := GetRegistrationDTO(t, mw)
	_, err := mw.HandleRegister(reg, LoginContext{IP: "127.0.0.1", Region: "SF"})
	require.Nil(t, err)

	login := GetLoginDTO(t, reg, mw)

	token, err := mw.HandleLogin(login, LoginContext{IP: "127.0.0.1", Region: "SF"})
	assert.Nil(t, err)
	assert.NotNil(t, token)

	reg.Username += "111"
	login = GetLoginDTO(t, reg, mw)
	_, err = mw.HandleLogin(login, LoginContext{})
	assert.Equal(t, errUserUnknown, err)
}

////////////////////////////////////////////////////////////////

func TestRegistration(t *testing.T) {
	DoTestRegistration(t, GetMockMiddleware())

	if testOnRealEnvironment {
		mw := GetMiddleware(t)
		defer mw.Close()
		DoTestRegistration(t, mw)
	}
}

func TestLogin(t *testing.T) {
	DoTestLogin(t, GetMockMiddleware())
	if testOnRealEnvironment {
		mw := GetMiddleware(t)
		defer mw.Close()
		DoTestLogin(t, mw)
	}
}
