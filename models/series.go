package models

import (
	"fmt"
	"github.com/garfunkel/go-tvdb"
	"time"
)

type Series struct {
	ID          uint64
	Matched     bool
	ImdbID      string
	Status      string
	SeriesID    string
	SeriesName  string
	Language    string
	LastUpdated string
	Episodes    map[string]*Episode
	LocalName   string
	LocalPath   string
	Poster      string
}

func (s *Series) mapInfo(series tvdb.Series) {
	s.Matched = true
	s.ID = series.ID
	s.ImdbID = series.ImdbID
	s.Status = series.Status
	s.SeriesName = series.SeriesName
	s.Language = series.Language
	s.LastUpdated = series.LastUpdated
	s.Poster = series.Poster
}

// Get basic information from tvdbcom by ImdbID
func (s *Series) FetchInfoByImdbID(imdbID string) {
	// TODO Get just directory instead of full path if LocalPath is absolute
	// seriesListTvDb, err := tvdb.GetSeries(s.LocalName)
	series, err := tvdb.GetSeriesByIMDBID(imdbID)
	if err != nil {
		s.Matched = false
	} else {
		s.mapInfo(*series)
	}
}

// Get basic information from tvdbcom by ID
func (s *Series) FetchInfoByID(id uint64) {
	// TODO Get just directory instead of full path if LocalPath is absolute
	// seriesListTvDb, err := tvdb.GetSeries(s.LocalName)
	series, err := tvdb.GetSeriesByID(id)
	if err != nil {
		s.Matched = false
	} else {
		s.mapInfo(*series)
	}
}

// Get basic information from tvdbcom
func (s *Series) FetchInfoByName(name string) {
	// TODO Get just directory instead of full path if LocalPath is absolute
	seriesListTvDb, err := tvdb.GetSeries(name)
	if err != nil || len(seriesListTvDb.Series) == 0 {
		s.Matched = false
		// fmt.Println("Could not match")
	} else {
		// series := *seriesListTvDb.Series[0]
		series, err := tvdb.GetSeriesByID(seriesListTvDb.Series[0].ID)
		if err != nil {
			s.Matched = false
		} else {
			s.mapInfo(*series)
		}
	}
}

// Get detailed information from tvdbcom
func (s *Series) FetchDetails() {
	if s.Matched == true {
		seriesListTvDb, err := tvdb.GetSeries(s.LocalName)
		if err != nil || len(seriesListTvDb.Series) == 0 {
			// fmt.Println("Could not match")
			return
		}
		series := *seriesListTvDb.Series[0]
		series.GetDetail()
		// s.Episodes = make(map[SeasonEpisode]*Episode)
		s.Episodes = make(map[string]*Episode)

		for _, seasonEpisodes := range series.Seasons {
			for _, episode := range seasonEpisodes {
				episodeSimple := Episode{}

				episodeSimple.ID = episode.ID
				episodeSimple.EpisodeName = episode.EpisodeName
				episodeSimple.EpisodeNumber = episode.EpisodeNumber
				episodeSimple.FirstAired = episode.FirstAired
				episodeSimple.ImdbID = episode.ImdbID
				episodeSimple.Language = episode.Language
				episodeSimple.SeasonNumber = episode.SeasonNumber
				episodeSimple.LastUpdated = episode.LastUpdated
				episodeSimple.SeasonID = episode.SeasonID
				episodeSimple.SeriesID = episode.SeriesID

				if episode.FirstAired == "" {
					// fmt.Println("Missing first aired.")
				} else {
					aired, err := time.Parse("2006-01-02", episode.FirstAired)
					if err != nil {
						fmt.Println("Could not parse first aired.", err)
					} else {
						if aired.Before(time.Now()) {
							episodeSimple.HasAired = true
							// fmt.Println(series.SeriesName, "Season", episode.SeasonNumber, "Episode", episode.EpisodeNumber, "aired", episode.FirstAired)
						} else {
							episodeSimple.HasAired = false
							// fmt.Println(series.SeriesName, "Season", episode.SeasonNumber, "Episode", episode.EpisodeNumber, "not yet aired, airing on", episode.FirstAired)
						}
					}
				}
				// s.Episodes[SeasonEpisode{episode.SeasonNumber, episode.EpisodeNumber}] = &episodeSimple
				s.Episodes[fmt.Sprintf("%d_%d", episode.SeasonNumber, episode.EpisodeNumber)] = &episodeSimple
			}
		}
	} else {
		// TODO Log Error
	}
}
