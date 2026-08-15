package repositories

import (
	"randomshit/database"
	"randomshit/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDTO struct {
	Id        int
	Firstname string
	Point     int
}


func CreateUser() (userDto UserDTO, err error) {

	var user *models.Users

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Event{
			ID:        uuid.New(),
			UserId:    user.ID,
			EventType: "CREATE USER",
			Version:   0,
			Amount:    &user.Point,
			Newpoint:  &user.Point,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	return UserDTO{
		Id: int(user.ID),
		Firstname: user.Name,
		Point: user.Point,
	}, err
}
