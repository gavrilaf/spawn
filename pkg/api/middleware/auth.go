package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/api/utils"
)

type Auth struct {
	*api.Bridge
}

func CreateAuth(bridge *api.Bridge) Auth {
	return Auth{Bridge: bridge}
}

func (self Auth) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := utils.JwtFromHeader(c, types.TokenLookup)
		if err != nil {
			utils.HandleAuthError(c, http.StatusBadRequest, err)
		}

		token, err := utils.ParseToken(tokenStr, func(id string) (interface{}, error) {
			cl, err := self.ReadModel.FindClient(id)
			if err != nil {
				return nil, err
			}
			return cl.Secret, nil
		})

		if err != nil {
			utils.HandleAuthError(c, http.StatusUnauthorized, err)
			return
		}

		claims := utils.ClaimsFromToken(token)
		session, err := self.ReadModel.GetSession(claims.SessionID())
		if err != nil {
			utils.HandleAuthError(c, http.StatusUnauthorized, types.ErrSessionNotFound)
			return
		}

		c.Set("session_id", session.ID)
		c.Set("client_id", session.ClientID)
		c.Set("user_id", session.UserID)
		c.Set("device_id", session.DeviceID)

		c.Next()
	}
}
