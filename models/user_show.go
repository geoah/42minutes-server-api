package models

type UserShow struct {
	UserID   string `json:"user_id" db:"user_id"`
	ShowID   int    `json:"show_id" db:"show_id"`
	Library  bool   `json:"library" db:"library"`
	Favorite bool   `json:"favorite" db:"favorite"`
	Path     string `json:"path" db:"path"`
}
