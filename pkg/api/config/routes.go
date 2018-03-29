package config

type Endpoint struct {
	Path   string
	Method string
}

const (
	gAuth     = "auth"
	gUser     = "user"
	gProfile  = "profile"
	gAccounts = "accounts"
)

var (
	eAuthRegister = Endpoint{"/register", "PUT"}
	eAuthLogin    = Endpoint{"/login", "POST"}
	eAuthRefresh  = Endpoint{"/refresh_token", "POST"}

	eUserState         = Endpoint{"/state", "GET"}
	eUserLogout        = Endpoint{"logout", "POST"}
	eUserDevices       = Endpoint{"/devices", "GET"}
	eUserDevicesDelete = Endpoint{"/devices/:id", "DELETE"}

	eProfileGet         = Endpoint{"/", "GET"}
	eProfileUpdCountry  = Endpoint{"/country", "POST"}
	eProfileUpdPersonal = Endpoint{"/personal", "POST"}

	eAccountsGet      = Endpoint{"/", "GET"}
	eAccountsState    = Endpoint{"/state/:id", "GET"}
	eAccountsRegister = Endpoint{"/register", "PUT"}
	eAccountsSuspend  = Endpoint{"/suspend/:id", "POST"}
	eAccountsResume   = Endpoint{"/resume/:id", "POST"}
)

/*type EndpointDesc struct {
	Path       string
	Method     string
	NeedDevice bool
	NeedEmail  bool
	MinScope   int
}

var apiConfig = map[string]map[int]EndpointDesc{
	"auth": map[int]EndpointDesc{
		eAuthRegister: {"register", "POST", false, false, 0},
	},
	"user": map[int]EndpointDesc{
		eUserState: {"/state", "GET", true, false, 0},
	},
}*/
