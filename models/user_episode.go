package models

type UserEpisode struct {
	UserID    int    `json:"user_id" db:"user_id"`
	ShowID    int    `json:"show_id" db:"show_id"`
	SeasonID  int    `json:"season" db:"season"`
	EpisodeID int    `json:"episode" db:"episode"`
	Path      string `json:"path" db:"path"`
}
