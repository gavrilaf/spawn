package config

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/middleware"
	"github.com/gavrilaf/spawn/pkg/api/types"

	"github.com/gavrilaf/spawn/pkg/api/account"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/profile"
	"github.com/gavrilaf/spawn/pkg/api/user"

	"github.com/gin-gonic/gin"
)

func ConfigureRouter(router gin.IRouter, bridge *api.Bridge) {

	authMiddleware := middleware.CreateAuthMiddleware(bridge)

	authAPI := auth.CreateApi(bridge)
	profileAPI := profile.CreateApi(bridge)
	userAPI := user.CreateApi(bridge)
	accountsApi := account.CreateApi(bridge)

	auth := router.Group(gAuth)
	{
		addHandler(auth, eAuthRegister, authAPI.SignUp)
		addHandler(auth, eAuthLogin, authAPI.SignIn)
		addHandler(auth, eAuthRefresh, authAPI.RefreshToken)
	}

	user := router.Group(gUser)
	user.Use(authMiddleware.MiddlewareFunc())
	{
		addHandler(user, eUserState, userAPI.GetState)
		addHandler(user, eUserLogout, userAPI.Logout)
		addHandler(user, eUserDevices, userAPI.GetDevices)
		addHandler(user, eUserDevicesDelete, userAPI.DeleteDevice)
	}

	profile := router.Group(gProfile)
	profile.Use(authMiddleware.MiddlewareFunc())
	{
		addHandler(profile, eProfileGet, profileAPI.GetUserProfile)
		addHandler(profile, eProfileUpdCountry, profileAPI.UpdateUserCountry)
		addHandler(profile, eProfileUpdPersonal, profileAPI.UpdateUserPersonalInfo)
	}

	accounts := router.Group(gAccounts)
	accounts.Use(authMiddleware.MiddlewareFunc())
	{
		addHandler(accounts, eAccountsGet, accountsApi.GetAccounts)
		addHandler(accounts, eAccountsState, accountsApi.GetAccountState)
		addHandler(accounts, eAccountsRegister, accountsApi.RegisterAccount)
		addHandler(accounts, eAccountsSuspend, accountsApi.SuspendAccount)
		addHandler(accounts, eAccountsResume, accountsApi.ResumeAccount)
	}
}

func addHandler(g *gin.RouterGroup, e types.Endpoint, f gin.HandlerFunc) {
	g.Handle(e.Method, e.Path, f)
}
