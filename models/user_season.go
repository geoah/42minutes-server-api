package models

type UserSeason struct {
	UserID   int    `json:"user_id" db:"user_id"`
	ShowID   int    `json:"show_id" db:"show_id"`
	SeasonID int    `json:"season_id" db:"season_id"`
	Path     string `json:"path" db:"path"`
}
