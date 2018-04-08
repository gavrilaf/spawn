package config

import "github.com/gavrilaf/spawn/pkg/api/defs"

const (
	gAuth     = "auth"
	gUser     = "user"
	gProfile  = "profile"
	gAccounts = "accounts"
)

var (
	eAuthRegister = defs.Endpoint{Path: "/register", Method: "PUT"}
	eAuthLogin    = defs.Endpoint{Path: "/login", Method: "POST"}
	eAuthRefresh  = defs.Endpoint{Path: "/refresh_token", Method: "POST"}

	eUserState                = defs.Endpoint{Path: "/state", Method: "GET"}
	eUserLogout               = defs.Endpoint{Path: "/logout", Method: "POST"}
	eUserDevices              = defs.Endpoint{Path: "/devices", Method: "GET"}
	eUserDevicesDelete        = defs.Endpoint{Path: "/devices/:id", Method: "DELETE"}
	eUserDeviceGetConfirmCode = defs.Endpoint{Path: "/devices/:id/code", Method: "GET"}
	eUserDeviceConfirm        = defs.Endpoint{Path: "/devices/confirm", Method: "POST"}

	eProfileGet         = defs.Endpoint{Path: "/", Method: "GET"}
	eProfileUpdCountry  = defs.Endpoint{Path: "/country", Method: "POST"}
	eProfileUpdPersonal = defs.Endpoint{Path: "/personal", Method: "POST"}

	eAccountsGet      = defs.Endpoint{Path: "/", Method: "GET"}
	eAccountsState    = defs.Endpoint{Path: "/state/:id", Method: "GET"}
	eAccountsRegister = defs.Endpoint{Path: "/register", Method: "PUT"}
	eAccountsSuspend  = defs.Endpoint{Path: "/suspend/:id", Method: "POST"}
	eAccountsResume   = defs.Endpoint{Path: "/resume/:id", Method: "POST"}
)
