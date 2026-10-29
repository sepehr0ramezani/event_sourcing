package main

import (
	"log"
	"os"
	"time"

	"randomshit/internal/applications/leaderboard/auth"
	handlers "randomshit/internal/applications/leaderboard/handler"
	"randomshit/internal/applications/leaderboard/middleware"
	"randomshit/internal/applications/leaderboard/repository"
	"randomshit/internal/applications/leaderboard/service"
	"randomshit/internal/infrastructure/database"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.StartDB()
	if err != nil {
		log.Fatalf("start database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewRepo(db)
	userService := service.NewUserService(userRepo)
	tokenManager, err := auth.NewJWTManager(os.Getenv("SECRET_KEY"), 24*time.Hour)
	if err != nil {
		log.Fatalf("configure JWT: %v", err)
	}
	userHandler := handlers.NewUserHandler(userService, tokenManager, os.Getenv("COOKIE_SECURE") == "true")

	router := gin.Default()
	router.POST("/signup", userHandler.SignUp)
	router.POST("/login", userHandler.Login)
	router.GET("/leaderboard", userHandler.Leaderboard)
	router.POST("/points", middleware.CheckCookie(tokenManager), userHandler.AddPoint)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("start HTTP server: %v", err)
	}
}
