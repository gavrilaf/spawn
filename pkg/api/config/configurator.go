package config

import (
	"github.com/gin-gonic/gin"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/middleware"

	"github.com/gavrilaf/spawn/pkg/api/account"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/profile"
	"github.com/gavrilaf/spawn/pkg/api/user"
)

func ConfigureEngine(engine *gin.Engine, bridge *api.Bridge) {

	authMiddleware := middleware.CreateAuth(bridge)
	accessMiddleware := middleware.CreateAccess(bridge, ApiDefaultAccess, ApiAccessConfig)

	authAPI := auth.CreateApi(bridge)
	profileAPI := profile.CreateApi(bridge)
	userAPI := user.CreateApi(bridge)
	accountsApi := account.CreateApi(bridge)

	auth := engine.Group(gAuth)
	{
		addHandler(auth, eAuthRegister, nil, authAPI.SignUp)
		addHandler(auth, eAuthLogin, nil, authAPI.SignIn)
		addHandler(auth, eAuthRefresh, nil, authAPI.RefreshToken)
	}

	user := engine.Group(gUser)
	user.Use(authMiddleware.MiddlewareFunc())
	{
		addHandler(user, eUserState, &accessMiddleware, userAPI.GetState)
		addHandler(user, eUserLogout, &accessMiddleware, userAPI.Logout)
		addHandler(user, eUserDevices, &accessMiddleware, userAPI.GetDevices)
		addHandler(user, eUserDevicesDelete, &accessMiddleware, userAPI.DeleteDevice)
	}

	profile := engine.Group(gProfile)
	profile.Use(authMiddleware.MiddlewareFunc())
	{
		addHandler(profile, eProfileGet, &accessMiddleware, profileAPI.GetUserProfile)
		addHandler(profile, eProfileUpdCountry, &accessMiddleware, profileAPI.UpdateUserCountry)
		addHandler(profile, eProfileUpdPersonal, &accessMiddleware, profileAPI.UpdateUserPersonalInfo)
	}

	accounts := engine.Group(gAccounts)
	accounts.Use(authMiddleware.MiddlewareFunc())
	{
		addHandler(accounts, eAccountsGet, &accessMiddleware, accountsApi.GetAccounts)
		addHandler(accounts, eAccountsState, &accessMiddleware, accountsApi.GetAccountState)
		addHandler(accounts, eAccountsRegister, &accessMiddleware, accountsApi.RegisterAccount)
		addHandler(accounts, eAccountsSuspend, &accessMiddleware, accountsApi.SuspendAccount)
		addHandler(accounts, eAccountsResume, &accessMiddleware, accountsApi.ResumeAccount)
	}
}

func addHandler(g *gin.RouterGroup, e defs.Endpoint, access *middleware.Access, f gin.HandlerFunc) {
	if access != nil {
		key := defs.GetEndpointKey(g.BasePath(), e)
		addRoute := func(c *gin.Context) {
			c.Set(defs.EndpointKey, key)
		}

		g.Handle(e.Method, e.Path, addRoute, access.MiddlewareFunc(), f)
	} else {
		g.Handle(e.Method, e.Path, f)
	}
}
