package config

import "github.com/gavrilaf/spawn/pkg/api/types"

const (
	gAuth     = "auth"
	gUser     = "user"
	gProfile  = "profile"
	gAccounts = "accounts"
)

var (
	eAuthRegister = types.Endpoint{Path: "/register", Method: "PUT"}
	eAuthLogin    = types.Endpoint{Path: "/login", Method: "POST"}
	eAuthRefresh  = types.Endpoint{Path: "/refresh_token", Method: "POST"}

	eUserState         = types.Endpoint{Path: "/state", Method: "GET"}
	eUserLogout        = types.Endpoint{Path: "logout", Method: "POST"}
	eUserDevices       = types.Endpoint{Path: "/devices", Method: "GET"}
	eUserDevicesDelete = types.Endpoint{Path: "/devices/:id", Method: "DELETE"}

	eProfileGet         = types.Endpoint{Path: "/", Method: "GET"}
	eProfileUpdCountry  = types.Endpoint{Path: "/country", Method: "POST"}
	eProfileUpdPersonal = types.Endpoint{Path: "/personal", Method: "POST"}

	eAccountsGet      = types.Endpoint{Path: "/", Method: "GET"}
	eAccountsState    = types.Endpoint{Path: "/state/:id", Method: "GET"}
	eAccountsRegister = types.Endpoint{Path: "/register", Method: "PUT"}
	eAccountsSuspend  = types.Endpoint{Path: "/suspend/:id", Method: "POST"}
	eAccountsResume   = types.Endpoint{Path: "/resume/:id", Method: "POST"}
)
