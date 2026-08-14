package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId    uint      `gorm:"not null;index"`
	EventType string    `gorm:"not null"`
	Version   int       `gorm:"not null"`
	CreatedAt time.Time
	Amount    *int
	Oldpoint  *int
	Newpoint  *int
}
