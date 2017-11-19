package auth

import (
	"testing"
	"time"

	//"github.com/stretchr/testify/assert"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	mdl "github.com/gavrilaf/spawn/pkg/dbx/model"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	tClientID = "client_test"
	tDeviveID = "device1"
	tPsw      = "password"
)

var storageMock = NewStorageMock(env.GetEnvironment("Test"))

func GetMiddleware() *Middleware {
	log := logrus.New()

	middleware := &Middleware{Timeout: time.Minute, MaxRefresh: time.Hour, Stg: storageMock, Log: log}

	return middleware
}

func GetClient(t *testing.T) mdl.Client {
	p, err := storageMock.FindClient(tClientID)
	require.Nil(t, err)
	return p
}

func GetRegistrationDTO(t *testing.T) RegisterDTO {
	client := GetClient(t)
	username := uuid.NewV4().String()
	sign := cryptx.GenerateSignature(client.ID+tDeviveID+username, client.Secret)
	return RegisterDTO{ClientID: client.ID, DeviceID: tDeviveID, Username: username, Password: tPsw, Signature: sign}
}

func GetLoginDTO(t *testing.T, reg RegisterDTO) LoginDTO {
	client := GetClient(t)
	sign := cryptx.GenerateSignature(client.ID+reg.DeviceID+reg.Username, client.Secret)

	return LoginDTO{ClientID: client.ID,
		DeviceID:  reg.DeviceID,
		AuthType:  AuthTypeSimple,
		Username:  reg.Username,
		Password:  reg.Password,
		Signature: sign}
}

////////////////////////////////////////////////////////////////

func TestRegistration(t *testing.T) {
	middleware := GetMiddleware()

	p := GetRegistrationDTO(t)

	_, err := middleware.HandleRegister(p)
	require.Nil(t, err)

	// Already registered
	_, err = middleware.HandleRegister(p)
	require.Equal(t, errUserAlreadyExist, err)

	// Invalid signature
	p.Signature += "111"
	_, err = middleware.HandleRegister(p)
	require.Equal(t, errInvalidSignature, err)
}

func TestLogin(t *testing.T) {
	middleware := GetMiddleware()

	reg := GetRegistrationDTO(t)
	_, err := middleware.HandleRegister(reg)
	require.Nil(t, err)

	login := GetLoginDTO(t, reg)

	token, err := middleware.HandleLogin(login)
	require.Nil(t, err)
	require.NotNil(t, token)

	reg.Username += "111"
	login = GetLoginDTO(t, reg)
	_, err = middleware.HandleLogin(login)
	require.Equal(t, errUserUnknown, err)
}
