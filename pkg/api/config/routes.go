package config

/*auth := router.Group("/auth")
{
	auth.POST("/register", authMiddleware.RegisterHandler)
	auth.POST("/login", authMiddleware.LoginHandler)
	auth.POST("/refresh_token", authMiddleware.RefreshHandler)
}

user := router.Group("user")
user.Use(authMiddleware.MiddlewareFunc())
{
	user.GET("/state", userAPI.GetState)
	user.POST("/logout", userAPI.Logout)
	user.GET("/devices", userAPI.GetDevices)
	user.DELETE("/devices/:id", userAPI.DeleteDevice)
}

profile := router.Group("profile")
profile.Use(authMiddleware.MiddlewareFunc())
{
	profile.GET("/", profileAPI.GetUserProfile)

	profile.POST("/country", profileAPI.UpdateUserCountry)
	profile.POST("/personal", profileAPI.UpdateUserPersonalInfo)
}

accounts := router.Group("accounts")
accounts.Use(authMiddleware.MiddlewareFunc())
{
	accounts.GET("/", accountsApi.GetAccounts)
	accounts.GET("/state/:id", accountsApi.GetAccountState)

	accounts.POST("/register", accountsApi.RegisterAccount)

	accounts.POST("/suspend/:id", accountsApi.SuspendAccount)
	accounts.POST("/resume/:id", accountsApi.ResumeAccount)
}
*/

const (
	gAuth         = "auth"
	eAuthRegister = "/register"
	eAuthLogin    = "/login"
	eAuthRefresh  = "/refresh_token"

	gUser              = "user"
	eUserState         = "/state"
	eUserLogout        = "logout"
	eUserDevices       = "/devices"
	eUserDevicesDelete = "/devices/:id"

	gProfile            = "profile"
	eProfileGet         = "/"
	eProfileUpdCountry  = "/country"
	eProfileUpdPersonal = "/personal"

	gAccounts         = "accounts"
	eAccountsGet      = "/"
	eAccountsState    = "/state/:id"
	eAccountsRegister = "/register"
	eAccountsSuspend  = "/suspend/:id"
	eAccountsResume   = "/resume/:id"
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
