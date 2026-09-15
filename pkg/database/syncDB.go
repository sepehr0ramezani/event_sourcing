package database

import "randomshit/pkg/models"

func SyncDateBase() {
	err := DB.AutoMigrate(
		&models.Users{},
		&models.Event{},
	)
	if err != nil {
		panic("faild to migirate")
	}
}
