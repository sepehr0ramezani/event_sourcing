package main

import (
	"randomshit/database"
	_ "randomshit/docs"
	"randomshit/handlers"
	"randomshit/middleware"

	files "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	database.StartDB()
	database.SyncDateBase()
}

// @title			preview of this shit
// @version		1.0
// @description	swagger for this web
// @host			localhost:8080
// @BasePath		/
func main() {

	router := gin.Default()

	router.GET("/leaderboard", middleware.CheckCookie, handlers.Leaderboard)
	router.POST("/signup", handlers.SignUp)
	router.POST("/login", handlers.Login)
	router.POST("/addpoint", middleware.CheckCookie, handlers.AddPoint)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	router.Run()
}
