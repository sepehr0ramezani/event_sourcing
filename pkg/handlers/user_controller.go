package handlers

import (
	"net/http"
	"randomshit/pkg/models"
	"randomshit/pkg/repositories"
	"randomshit/pkg/service"

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

	users, err := service.Leaderboard(c.Request.Context())
	if err != nil {
		c.JSON(400, gin.H{
			"error": "cant get leaderboard",
		})
		return
	}
	c.JSON(200, users)
}

func SignUp(c *gin.Context) {
	var input models.Body
	if c.Bind(&input) != nil {
		c.JSON(400, gin.H{
			"error": "faild to read request",
		})
		return
	}
	err, hashpass := service.HashPassword(input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "faild to hash password",
		})
		return
	}
	err = repositories.CreateUser(c, input, hashpass)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, input)
}

func Login(c *gin.Context) {
	var input models.Body
	if c.Bind(&input) != nil {
		c.JSON(400, gin.H{
			"error": "failed to read request",
		})
		return
	}
	err, token := repositories.Login(input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "can't find your user",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"token",
		token,
		86400,
		"/",
		"",
		false,
		true,
	)
}

func AddPoint(c *gin.Context) {
	var point struct {
		Point int `json:"point"`
	}
	if err := c.Bind(&point); err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	UserID, exist := c.Get("user_id")
	if !exist {
		c.JSON(401, gin.H{
			"error": "user not found",
		})
		return
	}
	err := repositories.Updatepoints(c, UserID, point.Point)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "can't update your points",
		})
		return
	}
	c.JSON(200, point)
}
