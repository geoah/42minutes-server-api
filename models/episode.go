package models

import (
	"github.com/hobeone/gotrakt"
)

type Episode struct {
	ShowID        int               `json:"show_id" db:"show_id"`
	Season        int               `json:"season" db:"season"`
	Episode       int               `json:"episode" db:"episode"`
	Number        int               `json:"number" db:"number"`
	TvdbID        int               `json:"tvdb_id" db:"tvdb_id"`
	Title         string            `json:"title" db:"title"`
	Overview      string            `json:"overview" db:"overview"`
	FirstAired    int               `json:"first_aired" db:"first_aired"`
	FirstAiredIso string            `json:"first_aired_iso" db:"first_aired_iso"`
	FirstAiredUtc int               `json:"first_aired_utc" db:"first_aired_utc"`
	URL           string            `json:"url" db:"url"`
	Screen        string            `json:"screen" db:"screen"`
	Images        map[string]string `json:"images" db:"-"`
}

func (e *Episode) MapInfo(episode gotrakt.Episode) {
	e.Season = episode.Season
	e.Episode = episode.Episode
	e.Number = episode.Number
	e.TvdbID = episode.TvdbID
	e.Title = episode.Title
	e.Overview = episode.Overview
	e.FirstAired = episode.FirstAired
	e.FirstAiredIso = episode.FirstAiredIso
	e.FirstAiredUtc = episode.FirstAiredUtc
	e.URL = episode.URL
	e.Screen = episode.Screen
	e.Images = episode.Images
}
