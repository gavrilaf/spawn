package profile

import (
	"github.com/gin-gonic/gin"
	//"github.com/sirupsen/logrus"
)

func (api *ProfileApiImpl) WhoAmI(c *gin.Context) {
	userId := c.GetString("user_id")
	clientId := c.GetString("client_id")

	c.JSON(200, gin.H{
		"user_id":   userId,
		"client_id": clientId,
	})
}

func (p *ProfileApiImpl) ConfirmDevice(c *gin.Context) {
	userId := c.GetString("user_id")
	clientId := c.GetString("client_id")

	c.JSON(200, gin.H{
		"user_id":   userId,
		"client_id": clientId,
	})
}
