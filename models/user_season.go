package models

type UserSeason struct {
	UserID   int    `json:"user_id" db:"user_id"`
	ShowID   int    `json:"show_id" db:"show_id"`
	SeasonID int    `json:"season" db:"season"`
	Path     string `json:"path" db:"path"`
}
