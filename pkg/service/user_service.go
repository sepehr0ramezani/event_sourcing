package service

import (
	"context"
	"randomshit/pkg/models"
	"randomshit/pkg/repositories"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var userDB []models.Users

func Leaderboard(c context.Context) ([]redis.Z, error) {
	return repositories.Leaderboard(c)

}

func HashPassword(input models.Body) (error, string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err, ""
	}
	return nil, string(hash)
}
