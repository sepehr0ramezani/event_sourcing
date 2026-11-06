package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"randomshit/internal/applications/leaderboard/models"
	"randomshit/internal/applications/leaderboard/port"

	"github.com/redis/go-redis/v9"
)

const (
	leaderboardCacheKey = "leaderboard:v1"
	leaderboardCacheTTL = 30 * time.Second
)

var _ port.UserRepo = (*Repository)(nil)

type Repository struct {
	db      *sql.DB
	queries *Queries
	redis   *redis.Client
}

func NewRepo(db *sql.DB, redisClient *redis.Client) *Repository {
	return &Repository{
		db:      db,
		queries: New(db),
		redis:   redisClient,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.Users) (*models.Users, error) {
	createdUser, err := r.queries.CreateUser(ctx, CreateUserParams{
		Username:     user.Username,
		HashPassword: user.Password,
	})
	if err != nil {
		return nil, err
	}
	if r.redis != nil {
		_ = r.redis.Del(ctx, leaderboardCacheKey).Err()
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	userinfo, err := qtx.AddPoint(ctx, AddPointParams{
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

	if err := qtx.AddESPoint(ctx, AddESPointParams{
		UserID: int32(id),
		Amount: int32(sumpoint),
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	if r.redis != nil {
		_ = r.redis.Del(ctx, leaderboardCacheKey).Err()
	}

	return &models.Users{
		Username: userinfo.Username,
		Point:    int(userinfo.Point.Int32),
	}, nil
}

func (r *Repository) Leaderboard(ctx context.Context) ([]models.LeaderboardEntry, error) {
	if r.redis != nil {
		cached, err := r.redis.Get(ctx, leaderboardCacheKey).Bytes()
		if err == nil {
			var entries []models.LeaderboardEntry
			if json.Unmarshal(cached, &entries) == nil {
				return entries, nil
			}
		}
		// Cache miss or Redis error: rebuild the result from Postgres.
	}

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

	if r.redis != nil {
		if data, err := json.Marshal(entries); err == nil {
			_ = r.redis.Set(ctx, leaderboardCacheKey, data, leaderboardCacheTTL).Err()
		}
	}

	return entries, nil
}
