package models

type UserShow struct {
	UserID int    `json:"user_id" db:"user_id"`
	ShowID int    `json:"show_id" db:"show_id"`
	Path   string `json:"path" db:"path"`
}
