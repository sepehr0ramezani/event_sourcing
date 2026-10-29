package models

type Users struct {
	ID       int
	Username string `json:"username"`
	Password string `json:"-"`
	Point    int    `json:"point"`
}

type LeaderboardEntry struct {
	Username string `json:"username"`
	Point    int    `json:"point"`
}

type Body struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
