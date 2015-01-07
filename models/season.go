package models

import (
	"github.com/hobeone/gotrakt"
)

type Season struct {
	ShowID   int       `json:"show_id" db:"show_id"`
	Season   int       `json:"season" db:"season"`
	URL      string    `json:"url" db:"url"`
	Poster   string    `json:"poster" db:"poster"`
	Episodes []Episode `json:"episodes" db:"-"`
}

func (s *Season) MapInfo(season gotrakt.Season) {
	s.Season = season.Season
	s.URL = season.URL
	s.Poster = season.Poster

	s.Episodes = make([]Episode, 0)
	for _, anepisode := range season.Episodes {
		episode := Episode{}
		episode.ShowID = s.ShowID
		episode.MapInfo(anepisode)
		s.Episodes = append(s.Episodes, episode)
	}
}
