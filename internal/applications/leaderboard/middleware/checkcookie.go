package middleware

import (
	"randomshit/internal/applications/leaderboard/auth"

	"github.com/gin-gonic/gin"
)

func CheckCookie(tokens *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "you aren't logged in"})
			return
		}

		userID, err := tokens.UserID(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
