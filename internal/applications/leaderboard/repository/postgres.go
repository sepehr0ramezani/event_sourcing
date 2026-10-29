package repository

import (
	"context"
	"database/sql"
	"errors"

	"randomshit/internal/applications/leaderboard/models"
	"randomshit/internal/applications/leaderboard/port"
)

var _ port.UserRepo = (*Repository)(nil)

type Repository struct {
	queries *Queries
}

func NewRepo(db DBTX) *Repository {
	return &Repository{queries: New(db)}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.Users) (*models.Users, error) {
	createdUser, err := r.queries.CreateUser(ctx, CreateUserParams{
		Username:     user.Username,
		HashPassword: user.Password,
	})
	if err != nil {
		return nil, err
	}

	return &models.Users{
		ID:       int(createdUser.ID),
		Username: createdUser.Username,
	}, nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*models.Users, error) {
	userinfo, err := r.queries.Login(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}

	return &models.Users{
		ID:       int(userinfo.ID),
		Username: userinfo.Username,
		Password: userinfo.HashPassword,
	}, nil
}

func (r *Repository) AddPoint(ctx context.Context, id, sumpoint int) (*models.Users, error) {
	userinfo, err := r.queries.AddPoint(ctx, AddPointParams{
		Point: sql.NullInt32{
			Int32: int32(sumpoint),
			Valid: true,
		},
		ID: int32(id),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}

	return &models.Users{
		Username: userinfo.Username,
		Point:    int(userinfo.Point.Int32),
	}, nil
}

func (r *Repository) Leaderboard(ctx context.Context) ([]models.LeaderboardEntry, error) {
	leaderBoard, err := r.queries.LeaderBoard(ctx)
	if err != nil {
		return nil, err

	}

	var entries []models.LeaderboardEntry
	for _, row := range leaderBoard {

		entries = append(entries, models.LeaderboardEntry{
			Username: row.Username,
			Point:    int(row.Point.Int32),
		})

	}

	return entries, nil
}
