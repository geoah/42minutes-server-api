package models

type UserFile struct {
	ID           string `json:"id" db:"id"`
	UserID       int    `json:"user_id" db:"user_id"`
	RelativePath string `json:"relative_path" db:"relative_path"`
	FullPathHash string `json:"full_path_hash" db:"full_path_hash"`
	Hash         string `json:"hash" db:"hash"`
	Processed    bool   `json:"processed" db:"processed"`
	Matched      bool   `json:"matched" db:"matched"`
	ShowID       int    `json:"show_id" db:"show_id"`
	SeasonID     int    `json:"season_id" db:"season_id"`
	EpisodeID    int    `json:"episode_id" db:"episode_id"`
}
