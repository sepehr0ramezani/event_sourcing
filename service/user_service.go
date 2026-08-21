package service

import (
	"randomshit/models"
	"randomshit/repositories"

	"golang.org/x/crypto/bcrypt"
)

var userDB []models.Users

func Leaderboard() (err error, users []models.Users) {
	err = repositories.Leaderboard(&userDB)
	if err != nil {
		return err, nil
	}
	return err, userDB
}

func HashPassword(input models.Body) (error, string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err, ""
	}
	return nil, string(hash)
}
