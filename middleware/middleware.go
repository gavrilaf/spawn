package middleware

import (
	"github.com/gin-gonic/gin"
)

func HandleLogin(p *Login) (string, error) {

	return "user-id-1", nil
}

func CheckAccess(userId string, c *gin.Context) bool {
	return true
}
