package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gavrilaf/spawn/pkg/api"
	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/api/utils"
)

type AuthMiddleware struct {
	*api.Bridge
}

func CreateAuthMiddleware(bridge *api.Bridge) AuthMiddleware {
	return AuthMiddleware{Bridge: bridge}
}

func (self AuthMiddleware) MiddlewareFunc() gin.HandlerFunc {
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

		// TODO: Add later
		/*if err = mw.checker.CheckAccess(session, c); err != nil {
			mw.handleError(c, http.StatusForbidden, err)
			return
		}*/

		c.Set("session_id", session.ID)
		c.Set("client_id", session.ClientID)
		c.Set("user_id", session.UserID)
		c.Set("device_id", session.DeviceID)

		c.Next()
	}
}
