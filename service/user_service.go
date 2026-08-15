package service

import (
	"randomshit/database"
	"randomshit/models"
	"randomshit/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Addpoint(user *models.Users) (*models.Users, error) {

	var olduser models.Users

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("name = ?", user.Name).First(&olduser).Error; err != nil {
			return err
		}

		if err := tx.Save(&models.Users{
			ID:    olduser.ID,
			Name:  olduser.Name,
			Point: user.Point + olduser.Point,
		}).Error; err != nil {
			return err
		}

		var lastVersion int

		if err := tx.Model(&models.Event{}).
			Where("user_id = ?", olduser.ID).
			Select("COALESCE(MAX(version),0)").Scan(&lastVersion).
			Error; err != nil {

			return err

		}

		lastVersion += 1

		newscore := olduser.Point + user.Point

		if err := tx.Create(&models.Event{

			ID:        uuid.New(),
			UserId:    olduser.ID,
			EventType: "ADD POINT",
			Version:   lastVersion,
			Amount:    &user.Point,
			Oldpoint:  &olduser.Point,
			Newpoint:  &newscore,
		}).Error; err != nil {
			return err
		}

		return nil

	})

	return &olduser, err

}

func CreateUser() (user_id int, err error) {
	user_dto, err := repositories.CreateUser()
	if err != nil {
		return  0, err
	}
	return user_dto.Id, err

}
