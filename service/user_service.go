package service

import (
	"randomshit/database"
	"randomshit/models"
	"randomshit/repositories"
)

var userDB []models.Users

func Leaderboard() (err error, users []models.Users) {
	err = database.DB.Order("point DESC").Find(&userDB).Error
	if err != nil {
		return err, nil
	}
	return err, userDB
}

func NewUpdate(input models.Users) (err error) {
	err = database.DB.Find(&userDB).Error
	if err != nil {
		return err
	}

	//trying to find this user in db
	for _, user := range userDB {
		if user.Name == input.Name {
			//if it's find it create user
			err = repositories.UpdateOldUser(input)
			if err != nil {
				return err
			}
			return nil
		}
	}
	// else update user in db
	err = repositories.CreateUser(input)
	if err != nil {
		return err
	}
	return nil
}
