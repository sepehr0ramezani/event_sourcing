package repositories

import (
	"errors"
	"os"
	"randomshit/database"
	"randomshit/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(user models.Body, hashpass string) error {

	newuser := models.Users{
		Username: user.Username,
		Password: hashpass,
		Point:    user.Point,
	}
	var wasindb models.Users
	result := database.DB.First(&wasindb, "username = ?", newuser.Username)

	if result.Error == nil {
		return errors.New("this user exsit")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&newuser).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Event{
			UUID:      uuid.New(),
			UserId:    uint(newuser.ID),
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

func Login(body models.Body) (error, string) {

	//find user
	var user models.Users
	database.DB.First(&user, "username = ?", body.Username)
	if user.ID == 0 {
		return errors.New("can't find it"), ""
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return err, ""
	}

	//jwt token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * 30 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sec_key := os.Getenv("SECRET_KEY")

	var tokenstring string
	tokenstring, err = token.SignedString([]byte(sec_key))
	if err != nil {
		return err, ""
	}

	return nil, tokenstring
}

func Leaderboard(userDB *[]models.Users) error {
	err := database.DB.Order("point DESC").Find(&userDB).Error
	if err != nil {
		return err
	}
	return nil
}

func Updatepoints(userid any, points int) error {
	id, exist := userid.(int)
	var userinfo models.Users
	if !exist {
		return errors.New("cant find user id")
	}

	result := database.DB.First(&userinfo, "id = ?", id)
	if result.Error != nil {
		return errors.New("its not on db")
	}
	userinfo.Point += points
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&userinfo).Error; err != nil {
			return err
		}
		var lastEvent models.Event
		result = tx.
			Where("user_id = ?", id).
			Order("version DESC").
			First(&lastEvent)

		if result.Error != nil {
			return result.Error
		}
		lastEvent.Oldpoint = lastEvent.Newpoint
		lastEvent.Amount = &points
		*lastEvent.Newpoint += points
		lastEvent.EventType = "UPDATE POINTS"
		lastEvent.Version += 1
		lastEvent.UUID = uuid.New()
		lastEvent.CreatedAt = time.Now()

		if err := tx.Create(&lastEvent).Error; err != nil {
			return err
		}
		return nil
	})
}
