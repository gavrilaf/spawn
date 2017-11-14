package api

import (
	//"github.com/gavrilaf/spawn/pkg/auth"
	"github.com/gin-gonic/gin"
	//"github.com/sirupsen/logrus"
)

type CustomerDTO struct {
	Username string
}

func (p *SpawnApi) WhoAmI(c *gin.Context) {
	userId := c.GetString("user_id")
	clientId := c.GetString("client_id")

	c.JSON(200, gin.H{
		"user_id":   userId,
		"client_id": clientId,
	})
}
