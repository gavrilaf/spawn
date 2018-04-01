package ginx

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"github.com/gavrilaf/spawn/pkg/api/defs"
)

type Claims struct {
	claims jwt.MapClaims
}

func ClaimsFromToken(token *jwt.Token) Claims {
	return Claims{claims: token.Claims.(jwt.MapClaims)}
}

func (c Claims) SessionID() string {
	return c.claims["session_id"].(string)
}

func (c Claims) ClientID() string {
	return c.claims["aud"].(string)
}

//////////////////////////////////////////////////////////////////////////////////////////

// ParseToken - validate and parse jwt token
// token -
// secretFunc - get client secret by client id stored in the token
func ParseToken(token string, secretFunc func(string) (interface{}, error)) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		claims := ClaimsFromToken(token)
		return secretFunc(claims.ClientID())
	})
}

func JwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", defs.ErrInvalidRequest
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == defs.TokenHeadName) {
		return "", defs.ErrInvalidRequest
	}

	return parts[1], nil
}
