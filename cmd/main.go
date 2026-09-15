package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"randomshit/pkg/database"
	_ "randomshit/pkg/docs"
	"randomshit/pkg/handlers"
	"randomshit/pkg/middleware"

	files "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	database.StartDB()
	database.SyncDateBase()
	database.ConnectRedis()
	err := database.PingRedis()
	if err != nil {
		log.Fatal("can't connect to redis")
		return
	}

}

// @title			preview of this shit
// @version		1.0
// @description	swagger for this web
// @host			localhost:8080
// @BasePath		/
func main() {

	router := gin.Default()

	go http.ListenAndServe(":6060", nil)
	router.GET("/leaderboard", middleware.CheckCookie, handlers.Leaderboard)
	router.POST("/signup", handlers.SignUp)
	router.POST("/login", handlers.Login)
	router.POST("/addpoint", middleware.CheckCookie, handlers.AddPoint)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	router.Run()
}
