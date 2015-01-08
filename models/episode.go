package models

import (
	"github.com/hobeone/gotrakt"
)

type Episode struct {
	ShowID        int               `json:"show_id" db:"show_id"`
	Season        int               `json:"season" db:"season"`
	Episode       int               `json:"episode" db:"episode"`
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

func (e *Episode) MapInfo(traktEpisode gotrakt.Episode) {
	e.Season = traktEpisode.Season
	e.Episode = traktEpisode.Number
	e.TvdbID = traktEpisode.TvdbID
	e.Title = traktEpisode.Title
	e.Overview = traktEpisode.Overview
	e.FirstAired = traktEpisode.FirstAired
	e.FirstAiredIso = traktEpisode.FirstAiredIso
	e.FirstAiredUtc = traktEpisode.FirstAiredUtc
	e.URL = traktEpisode.URL
	e.Screen = traktEpisode.Screen
	e.Images = traktEpisode.Images
}
