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
	router.StaticFile("/leaderboard", "./frontend/leadboard.html")
	router.GET("/api/leaderboard", handlers.Leaderboard)
	router.POST("/crud", handlers.Does_it_exsit)
	router.Run()
}
