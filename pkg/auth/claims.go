package auth

import (
	"gopkg.in/dgrijalva/jwt-go.v3"
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
