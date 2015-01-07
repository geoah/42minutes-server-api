package models

type UserSeason struct {
	UserID   int    `json:"user_id" db:"user_id"`
	SeasonID int    `json:"season_id" db:"season_id"`
	Path     string `json:"path" db:"path"`
}
