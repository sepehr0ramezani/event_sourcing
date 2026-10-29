package port

import (
	"context"
	"errors"

	"randomshit/internal/applications/leaderboard/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrInvalidPoints = errors.New("points must be positive")

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.Users) (*models.Users, error)
	FindByUsername(ctx context.Context, username string) (*models.Users, error)
	AddPoint(ctx context.Context, id, sumpoint int) (*models.Users, error)
	Leaderboard(ctx context.Context) ([]models.LeaderboardEntry, error)
}
