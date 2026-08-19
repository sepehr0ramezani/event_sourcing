package handlers

import (
	"randomshit/models"
	"randomshit/service"

	"github.com/gin-gonic/gin"
)

// @Summary		leaderboard of users
// @Description	order users by their points
// @Tags			leaderboard
// @Produce		json
// @Success		200	{array}		[]models.Users		"Users ordered by points from highest to lowest"
// @failure		400	{object}	map[string]string	"error"
// @Router			/leaderboard [get]
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

// @Summary		update in db
// @Description	update or create a user
// @Tags			Updateuser
// @Produce		json
// @Param			users	body		models.Users		true	"The user you want to create or give points to"
// @Success		200		{object}	models.Users		"updated user"
// @failure		400		{object}	map[string]string	"error"
// @Router			/crud [post]
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
