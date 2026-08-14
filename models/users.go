package models

type Users struct {
	ID    uint
	Name  string `json:"name" binding:"required"`
	Point int    `json:"point" binding:"required"`
}
