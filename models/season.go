package models

import (
	"github.com/hobeone/gotrakt"
	"log"
)

type Season struct {
	ShowID   int       `json:"show_id" db:"show_id"`
	Season   int       `json:"season" db:"season"`
	URL      string    `json:"url" db:"url"`
	Poster   string    `json:"poster" db:"poster"`
	Episodes []Episode `json:"episodes" db:"-"`
}

func (s *Season) MapInfo(traktSeason gotrakt.Season) {
	s.Season = traktSeason.Season
	s.URL = traktSeason.URL
	s.Poster = traktSeason.Poster

	s.Episodes = make([]Episode, 0)
	for _, anepisode := range traktSeason.Episodes {
		episode := Episode{}
		episode.ShowID = s.ShowID
		episode.MapInfo(anepisode)
		log.Println(episode)
		s.Episodes = append(s.Episodes, episode)
	}
}
