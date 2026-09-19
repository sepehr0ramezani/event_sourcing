package repositories

import (
	"context"
	"errors"
	"os"
	"randomshit/pkg/database"
	"randomshit/pkg/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, user models.Body, hashpass string) error {

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

	if erro := database.RedisClient.ZAdd(ctx, "leaderboard", redis.Z{
		Score:  float64(user.Point),
		Member: user.Username,
	}).Err(); err != nil {
		return erro
	}

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

func Leaderboard(ctx context.Context) ([]redis.Z, error) {
	result, err := database.RedisClient.ZRevRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Updatepoints(ctx context.Context, userid any, points int) error {
	id, exist := userid.(int)
	var userinfo models.Users
	if !exist {
		return errors.New("its not a number")
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.First(&userinfo, "id = ?", id)
		if result.Error != nil {
			return errors.New("its not on db")
		}
		userinfo.Point += points
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
		newpoint := *lastEvent.Newpoint + points
		newEvent := models.Event{
			UUID:      uuid.New(),
			UserId:    lastEvent.UserId,
			EventType: "UPDATE POINTS",
			Version:   lastEvent.Version + 1,
			CreatedAt: time.Now(),
			Amount:    &points,
			Oldpoint:  &userinfo.Point,
			Newpoint:  &newpoint,
		}

		if err := tx.Create(&newEvent).Error; err != nil {
			return err
		}
		return nil
	})
	if err := database.RedisClient.ZAdd(ctx, "leaderboard", redis.Z{
		Score:  float64(userinfo.Point),
		Member: userinfo.Username,
	}).Err(); err != nil {
		return err
	}

	return err
}
