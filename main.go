package main

import (
	"randomshit/database"
	_ "randomshit/docs"
	"randomshit/handlers"

	files "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	database.StartDB()
}

// @title My API
// @version 1.0
// @description My first Gin API with Swagger
// @host localhost:8080
// @BasePath /
func main() {

	router := gin.Default()

	router.GET("/leaderboard", handlers.Leaderboard)
	router.POST("/crud", handlers.UpdateUsers)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	router.Run()
}
