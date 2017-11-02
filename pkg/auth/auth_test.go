package auth

import (
	"testing"
	"time"

	//"github.com/stretchr/testify/assert"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/gavrilaf/spawn/pkg/cryptx"
)

const (
	tClientID = "client_test"
	tDeviveID = "device1"
	tUsername = "user1"
	tPsw      = "password"
)

func GetMiddleware() *Middleware {
	log := logrus.New()
	storage := StorageFacade{Clients: NewClientsStorageMock(), Users: NewUsersStorageMock(), Sessions: NewMemorySessionsStorage()}
	middleware := &Middleware{Timeout: time.Minute, MaxRefresh: time.Hour, Storage: storage, Log: log}

	return middleware
}

func GetClient(t *testing.T) *Client {
	clients := NewClientsStorageMock()
	p, err := clients.FindClientByID(tClientID)
	require.Nil(t, err)
	return p
}

func GetRegistrationParcel(t *testing.T) *RegisterParcel {
	client := GetClient(t)
	sign := cryptx.GenerateSignature(client.ID()+tDeviveID+tUsername, client.Secret())
	return &RegisterParcel{ClientID: client.ID(), DeviceID: tDeviveID, Username: tUsername, Password: tPsw, Signature: sign}
}

func GetLoginParcel(t *testing.T, username string) *LoginParcel {
	if username == "" {
		username = tUsername
	}

	client := GetClient(t)
	sign := cryptx.GenerateSignature(client.ID()+tDeviveID+username, client.Secret())

	return &LoginParcel{ClientID: client.ID(), DeviceID: tDeviveID, AuthType: AuthTypeSimple, Username: username, Password: tPsw, Signature: sign}
}

////////////////////////////////////////////////////////////////

func TestRegistration(t *testing.T) {
	middleware := GetMiddleware()

	p := GetRegistrationParcel(t)

	err := middleware.HandleRegister(p)
	require.Nil(t, err)

	// Already registered
	err = middleware.HandleRegister(p)
	require.Equal(t, errUserAlreadyExist, err)

	// Invalid signature
	p.Signature += "111"
	err = middleware.HandleRegister(p)
	require.Equal(t, errInvalidSignature, err)
}

func TestLogin(t *testing.T) {
	middleware := GetMiddleware()

	reg := GetRegistrationParcel(t)
	err := middleware.HandleRegister(reg)
	require.Nil(t, err)

	login := GetLoginParcel(t, "")

	token, err := middleware.HandleLogin(login)
	require.Nil(t, err)
	require.NotNil(t, token)

	login = GetLoginParcel(t, tUsername+"12")
	_, err = middleware.HandleLogin(login)
	require.Equal(t, errUserUnknown, err)
}
