package main

import (
	"randomshit/database"
	"randomshit/handlers"

	"github.com/gin-gonic/gin"
)

func init() {
	database.StartDB()
}

func main() {
	router := gin.Default()
	router.GET("/leaderboard", handlers.Leaderboard)
	router.POST("/crud", handlers.UpdateUsers)
	router.Run()
}
