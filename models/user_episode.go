package models

type UserEpisode struct {
	UserID    string `json:"user_id" db:"user_id"`
	ShowID    int    `json:"show_id" db:"show_id"`
	Season    int    `json:"season" db:"season"`
	Episode   int    `json:"episode" db:"episode"`
	Available bool   `json:"available" db:"available"`
	Watched   bool   `json:"watched" db:"watched"`
	Path      string `json:"path" db:"path"`
}
