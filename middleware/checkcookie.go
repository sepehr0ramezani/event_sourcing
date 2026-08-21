package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CheckCookie(c *gin.Context) {
	tokenString, err := c.Cookie("token")

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "you aren't logged in",
		})
		return
	}

	jwtSecret := os.Getenv("SECRET_KEY")
	//???????????????????????????????????????????????????????????????????????????????????????????????????????
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(jwtSecret), nil
		},
	)

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "invalid token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "invalid claims",
		})
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)

	if !ok {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "invalid user id",
		})
		return
	}

	userID := int(userIDFloat)

	c.Set("user_id", userID)

	c.Next()
}
