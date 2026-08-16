package repositories

import (
	"randomshit/database"
	"randomshit/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(user models.Users) (err error) {

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

	return err
}

func UpdateOldUser(user models.Users) (err error) {
	var olduser models.Users

	err = database.DB.Transaction(func(tx *gorm.DB) error {

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

	return err
}
