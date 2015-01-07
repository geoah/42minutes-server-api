package models

type UserEpisode struct {
	UserID    int    `json:"user_id" db:"user_id"`
	EpisodeID int    `json:"episode_id" db:"episode_id"`
	Path      string `json:"path" db:"path"`
}
