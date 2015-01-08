package models

type UserFile struct {
	ID           string `json:"id" db:"id"`
	UserID       string `json:"user_id" db:"user_id"`
	RelativePath string `json:"relative_path" db:"relative_path"`
	FullPathHash string `json:"full_path_hash" db:"full_path_hash"`
	FileMd5      string `json:"file_md5" db:"file_md5"`
	Processed    bool   `json:"processed" db:"processed"`
	Matched      bool   `json:"matched" db:"matched"`
	ShowID       int    `json:"show_id" db:"show_id"`
	SeasonID     int    `json:"season_id" db:"season_id"`
	EpisodeID    int    `json:"episode_id" db:"episode_id"`
}
