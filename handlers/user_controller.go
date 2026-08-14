package handlers

import (
	"randomshit/database"
	"randomshit/models"
	"randomshit/service"

	"github.com/gin-gonic/gin"
)

var userDB []models.Users

func Leaderboard(c *gin.Context) {
	err := database.DB.Order("point DESC").Find(&userDB).Error
	if err != nil {
		panic("database get fucked!!!")
	}
	c.JSON(200, userDB)
}

func Does_it_exsit(c *gin.Context) {
	var input models.Users

	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = database.DB.Find(&userDB).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var was_in_db bool
	for _, value := range userDB {

		if value.Name == input.Name {
			was_in_db = true
			break
		}
	}
	//create user
	if !was_in_db {
		newuser := models.Users{
			Name:  input.Name,
			Point: input.Point,
		}
		var updateuser *models.Users
		updateuser, err = service.CreateUser(&newuser)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(200, &updateuser)
		return

	} else { //update user points
		newuser := models.Users{
			Name:  input.Name,
			Point: input.Point,
		}
		var updateuser *models.Users
		updateuser, err = service.Addpoint(&newuser)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, &updateuser)
	}

}
