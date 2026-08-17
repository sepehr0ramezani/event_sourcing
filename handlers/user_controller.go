package handlers

import (
	"randomshit/models"
	"randomshit/service"

	"github.com/gin-gonic/gin"
)

// @Summary Get leaderboard
// @Description Get the current leaderboard
// @Tags leaderboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /leaderboard [get]
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
