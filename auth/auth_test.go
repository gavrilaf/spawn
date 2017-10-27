package auth

import (
	//"fmt"
	"testing"
	"time"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gavrilaf/go-auth/auth/cerr"
	"github.com/gavrilaf/go-auth/auth/storage"
	"github.com/gavrilaf/go-auth/cryptx"
)

const (
	tClientID = "client_test"
	tDeviveID = "device1"
	tUsername = "user1"
	tPsw      = "password"
)

func GetMiddleware() *Middleware {
	storage := storage.StorageFacade{Clients: storage.NewClientsStorageMock(), Users: storage.NewUsersStorageMock(), Sessions: storage.NewMemorySessionsStorage()}
	middleware := &Middleware{Timeout: time.Minute, MaxRefresh: time.Hour, Storage: storage}

	return middleware
}

func GetClient(t *testing.T) *storage.Client {
	clients := storage.ClientsStorageMock{}
	p, err := clients.FindClientByID(tClientID)
	require.Nil(t, err)
	return p
}

func GetRegistrationParcel(t *testing.T) *RegisterParcel {
	client := GetClient(t)
	sign := cryptx.GenerateSignature(client.ID()+tDeviveID+tUsername, client.Secret())
	return &RegisterParcel{ClientID: client.ID(), DeviceID: tDeviveID, Username: tUsername, Password: tPsw, Signature: sign}
}

func GetLoginParcel(t *testing.T) *LoginParcel {
	client := GetClient(t)
	sign := cryptx.GenerateSignature(client.ID()+tDeviveID+tUsername, client.Secret())

	hashedPsw, err := cryptx.GenerateHashedPassword(tPsw)
	require.Nil(t, err)

	return &LoginParcel{ClientID: client.ID(), DeviceID: tDeviveID, Username: tUsername, SignedPassword: hashedPsw, Signature: sign}
}

////////////////////////////////////////////////////////////////

func TestRegistration(t *testing.T) {
	middleware := GetMiddleware()

	p := GetRegistrationParcel(t)

	err := middleware.HandleRegister(p)
	require.Nil(t, err)

	// Already registered
	err = middleware.HandleRegister(p)
	require.Equal(t, cerr.UserAlreadyExist, err)

	// Invalid signature
	p.Signature += "111"
	err = middleware.HandleRegister(p)
	require.Equal(t, cerr.InvalidSignature, err)
}

func TestLogin(t *testing.T) {
	middleware := GetMiddleware()

	reg := GetRegistrationParcel(t)
	err := middleware.HandleRegister(reg)
	require.Nil(t, err)

	login := GetLoginParcel(t)

	token, err := middleware.HandleLogin(login)
	require.Nil(t, err)
	require.NotNil(t, token)

	login.Username = login.Username + "12"
	_, err = middleware.HandleLogin(login)
	require.Equal(t, cerr.UserUnknown, err)
}
