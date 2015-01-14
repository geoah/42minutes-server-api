package models

type UserEpisode struct {
	UserID    string `json:"user_id" db:"user_id"`
	ShowID    int    `json:"show_id" db:"show_id"`
	SeasonID  int    `json:"season" db:"season"`
	EpisodeID int    `json:"episode" db:"episode"`
	Watched   bool   `json:"watched" db:"watched"`
	Path      string `json:"path" db:"path"`
}
