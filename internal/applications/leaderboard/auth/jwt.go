package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTManager creates and verifies the application's login tokens.
type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) (*JWTManager, error) {
	if secret == "" {
		return nil, errors.New("SECRET_KEY is not set")
	}
	if ttl <= 0 {
		return nil, errors.New("token lifetime must be positive")
	}

	return &JWTManager{secret: []byte(secret), ttl: ttl}, nil
}

type claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Create(userID int) (string, error) {
	if userID <= 0 {
		return "", errors.New("user ID must be positive")
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	})

	return token.SignedString(m.secret)
}

func (m *JWTManager) UserID(tokenString string) (int, error) {
	tokenClaims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid || tokenClaims.ExpiresAt == nil || tokenClaims.UserID <= 0 {
		return 0, errors.New("invalid token claims")
	}

	return tokenClaims.UserID, nil
}
