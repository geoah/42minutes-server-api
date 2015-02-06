package models

import "time"

type EpisodeRss struct {
	ShowTitle  string     `db:"show_title"`
	Title      string     `db:"title"`
	Season     int        `db:"season"`
	Episode    int        `db:"episode"`
	FirstAired *time.Time `db:"first_aired"`
	InfohashHd string     `db:"infohash_hd720p"`
	InfohashSd string     `db:"infohash_sd480p"`
}
