package handlers

import (
	"randomshit/models"
	"randomshit/service"

	"github.com/gin-gonic/gin"
)

func Leaderboard(c *gin.Context) {

	err, users := service.Leaderboard()
	if err != nil {
		c.JSON(400, gin.H{
			"error": "cant order users",
		})
		return
	}
	c.JSON(200, users)
}

func UpdateUsers(c *gin.Context) {
	var input models.Users

	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = service.NewUpdate(input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, input)

}
