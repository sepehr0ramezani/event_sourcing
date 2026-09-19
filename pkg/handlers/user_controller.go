package handlers

import (
	"net/http"
	"randomshit/pkg/models"
	"randomshit/pkg/repositories"
	"randomshit/pkg/service"

	"github.com/gin-gonic/gin"
)

// @Summary      Get leaderboard
// @Description  Get users ordered by their points from highest to lowest.
// @Tags         leaderboard
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} LeaderboardEntry "Leaderboard ordered by points"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /leaderboard [get]
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

// @Summary      Sign up
// @Description  Create a new user account. Initial points are optional and default to 0.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body models.Body true "Signup request"
// @Success      200 {object} models.Body "User created successfully"
// @Failure      400 {object} map[string]string "Invalid request or failed to create user"
// @Router       /signup [post]
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

// @Summary      Login
// @Description  Authenticate a user and set a JWT token in an HttpOnly cookie.
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} map[string]string "Login successful"
// @Failure      400 {object} map[string]string "Invalid credentials"
// @Router       /login [post]
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

// @Summary      Add points
// @Description  Add points to the authenticated user's account.
// @Tags         points
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body AddPointRequest true "Points to add"
// @Success      200 {object} AddPointRequest "Points added successfully"
// @Failure      400 {object} map[string]string "Failed to add points"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /points [post]
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
