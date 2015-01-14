package models

type UserSeason struct {
	UserID string `json:"user_id" db:"user_id"`
	ShowID int    `json:"show_id" db:"show_id"`
	Season int    `json:"season" db:"season"`
	Path   string `json:"path" db:"path"`
}
