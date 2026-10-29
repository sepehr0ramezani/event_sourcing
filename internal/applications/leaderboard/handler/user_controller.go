package handlers

import (
	"context"
	"errors"
	"net/http"

	"randomshit/internal/applications/leaderboard/models"
	"randomshit/internal/applications/leaderboard/port"

	"github.com/gin-gonic/gin"
)

type UserCreator interface {
	CreateUser(ctx context.Context, username, password string) (*models.Users, error)
	Login(ctx context.Context, username, password string) (*models.Users, error)
	Leaderboard(ctx context.Context) ([]models.LeaderboardEntry, error)
	AddPoint(ctx context.Context, userID int, points int) (int, error)
}

type UserHandler struct {
	userService  UserCreator
	tokenIssuer  TokenIssuer
	cookieSecure bool
}

type TokenIssuer interface {
	Create(userID int) (string, error)
}

func NewUserHandler(userService UserCreator, tokenIssuer TokenIssuer, cookieSecure bool) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenIssuer:  tokenIssuer,
		cookieSecure: cookieSecure,
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var input models.Body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"username": user.Username,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input models.Body

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.userService.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		if errors.Is(err, port.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not log in"})
		return
	}

	token, err := h.tokenIssuer.Create(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create login token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 24*60*60, "/", "", h.cookieSecure, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (h *UserHandler) Leaderboard(c *gin.Context) {
	entries, err := h.userService.Leaderboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func (h *UserHandler) AddPoint(c *gin.Context) {
	var input struct {
		Points int `json:"points" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "points must be a positive integer",
		})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})
		return
	}
	result, err := h.userService.AddPoint(c.Request.Context(), userID, input.Points)

	if err != nil {
		if errors.Is(err, port.ErrInvalidPoints) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "points must be a positive integer"})
			return
		}
		if errors.Is(err, port.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add points"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"point": result})

}
