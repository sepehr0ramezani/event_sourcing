package models

type Users struct {
	ID       int
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Point    int    `json:"point" binding:"required"`
}
