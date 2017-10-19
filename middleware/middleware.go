package middleware

import (
	types "github.com/gavrilaf/go-auth"
	"github.com/gin-gonic/gin"
)

func HandleLogin(p *types.Login) (string, error) {
	return "user-id-1", nil
}

func CheckAccess(userId string, c *gin.Context) bool {
	return true
}
