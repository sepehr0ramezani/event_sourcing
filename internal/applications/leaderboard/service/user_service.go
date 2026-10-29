package service

import (
	"context"
	"errors"

	"randomshit/internal/applications/leaderboard/models"
	"randomshit/internal/applications/leaderboard/port"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo port.UserRepo
}

func NewUserService(repo port.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, username, password string) (*models.Users, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &models.Users{
		Username: username,
		Password: string(passwordHash),
	}

	return s.repo.CreateUser(ctx, newUser)
}

func (s *UserService) Login(ctx context.Context, username, password string) (*models.Users, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, port.ErrUserNotFound) {
			return nil, port.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, port.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) Leaderboard(ctx context.Context) ([]models.LeaderboardEntry, error) {
	return s.repo.Leaderboard(ctx)
}

func (s *UserService) AddPoint(ctx context.Context, userID int, points int) (int, error) {
	if userID <= 0 {
		return 0, port.ErrUserNotFound
	}
	if points <= 0 || int64(points) > int64(1<<31-1) {
		return 0, port.ErrInvalidPoints
	}

	user, err := s.repo.AddPoint(ctx, userID, points)
	if err != nil {
		if errors.Is(err, port.ErrUserNotFound) {
			return 0, port.ErrUserNotFound
		}
		return 0, err
	}
	return user.Point, nil
}
