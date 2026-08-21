package models

type Body struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Point    int    `json:"point"`
}
