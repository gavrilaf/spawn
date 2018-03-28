package config

import (
	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/account"
	"github.com/gavrilaf/spawn/pkg/api/auth"
	"github.com/gavrilaf/spawn/pkg/api/profile"
	"github.com/gavrilaf/spawn/pkg/api/user"

	"github.com/gin-gonic/gin"
)

func ConfigureRouter(router gin.IRouter, bridge *api.Bridge) {
	authMiddleware := auth.CreateMiddleware(bridge)

	profileAPI := profile.CreateApi(bridge)
	userAPI := user.CreateApi(bridge)
	accountsApi := account.CreateApi(bridge)

	auth := router.Group(gAuth)
	{
		auth.POST(eAuthRegister, authMiddleware.RegisterHandler)
		auth.POST(eAuthLogin, authMiddleware.LoginHandler)
		auth.POST(eAuthRefresh, authMiddleware.RefreshHandler)
	}

	user := router.Group(gUser)
	user.Use(authMiddleware.MiddlewareFunc())
	{
		user.GET(eUserState, userAPI.GetState)
		user.POST(eUserLogout, userAPI.Logout)
		user.GET(eUserDevices, userAPI.GetDevices)
		user.DELETE(eUserDevicesDelete, userAPI.DeleteDevice)
	}

	profile := router.Group(gProfile)
	profile.Use(authMiddleware.MiddlewareFunc())
	{
		profile.GET(eProfileGet, profileAPI.GetUserProfile)
		profile.POST(eProfileUpdCountry, profileAPI.UpdateUserCountry)
		profile.POST(eProfileUpdPersonal, profileAPI.UpdateUserPersonalInfo)
	}

	accounts := router.Group(gAccounts)
	accounts.Use(authMiddleware.MiddlewareFunc())
	{
		accounts.GET(eAccountsGet, accountsApi.GetAccounts)
		accounts.GET(eAccountsState, accountsApi.GetAccountState)
		accounts.POST(eAccountsRegister, accountsApi.RegisterAccount)
		accounts.POST(eAccountsSuspend, accountsApi.SuspendAccount)
		accounts.POST(eAccountsResume, accountsApi.ResumeAccount)
	}
}
