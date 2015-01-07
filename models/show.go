package models

import (
	// "fmt"
	"github.com/garfunkel/go-tvdb"
	// "time"
)

type Show struct {
	ID            uint64            `json:"id" db:"id"`
	Title         string            `json:"title" db:"title"`
	Year          int               `json:"year" db:"year"`
	URL           string            `json:"url" db:"url"`
	FirstAired    int               `json:"first_aired" db:"first_aired"`
	Country       string            `json:"country" db:"country"`
	Overview      string            `json:"overview" db:"overview"`
	Runtime       int               `json:"runtime" db:"runtime"`
	Network       string            `json:"network" db:"network"`
	AirDay        string            `json:"air_day" db:"air_day"`
	AirTime       string            `json:"air_time" db:"air_time"`
	Certification string            `json:"certification" db:"certification"`
	ImdbID        string            `json:"imdb_id" db:"imdb_id"`
	TvdbID        uint64            `json:"tvdb_id" db:"tvdb_id"`
	TvrageID      int               `json:"tvrage_id" db:"tvrage_id"`
	Ended         bool              `json:"ended" db:"ended"`
	Images        map[string]string `json:"images" db:"images"`
	Genres        []string          `json:"genres" db:"genres"`
	Seasons       []Season          `json:"seasons" db:"-"`
	Matched       bool              `db:"-"`
}

func (s *Show) mapInfo(series tvdb.Series) {
	s.Matched = true
	s.ID = series.ID
	s.TvdbID = series.ID
	s.ImdbID = series.ImdbID
	s.Title = series.SeriesName
}

// Get basic information from tvdbcom by ImdbID
func (s *Show) FetchInfoByImdbID(imdbID string) {
	// TODO Get just directory instead of full path if LocalPath is absolute
	// ShowListTvDb, err := tvdb.GetShow(s.LocalName)
	series, err := tvdb.GetSeriesByIMDBID(imdbID)
	if err != nil {
		s.Matched = false
	} else {
		s.mapInfo(*series)
	}
}

// Get basic information from tvdbcom by ID
func (s *Show) FetchInfoByID(id uint64) {
	// TODO Get just directory instead of full path if LocalPath is absolute
	// ShowListTvDb, err := tvdb.GetShow(s.LocalName)
	series, err := tvdb.GetSeriesByID(id)
	if err != nil {
		s.Matched = false
	} else {
		s.mapInfo(*series)
	}
}

// Get basic information from tvdbcom
func (s *Show) FetchInfoByName(name string) {
	// TODO Get just directory instead of full path if LocalPath is absolute
	SeriesListTvDb, err := tvdb.GetSeries(name)
	if err != nil || len(SeriesListTvDb.Series) == 0 {
		s.Matched = false
		// fmt.Println("Could not match")
	} else {
		// Show := *SeriesListTvDb.Series[0]
		series, err := tvdb.GetSeriesByID(SeriesListTvDb.Series[0].ID)
		if err != nil {
			s.Matched = false
		} else {
			s.mapInfo(*series)
		}
	}
}

// Get detailed information from tvdbcom
// func (s *Show) FetchDetails() {
// 	if s.Matched == true {
// 		SeriesListTvDb, err := tvdb.GetSeries(s.Title)
// 		if err != nil || len(SeriesListTvDb.Series) == 0 {
// 			// fmt.Println("Could not match")
// 			return
// 		}
// 		series := *SeriesListTvDb.Series[0]
// 		series.GetDetail()
// 		// s.Episodes = make(map[SeasonEpisode]*Episode)
// 		s.Episodes = make(map[string]*Episode)

// 		for _, seasonEpisodes := range series.Seasons {
// 			for _, episode := range seasonEpisodes {
// 				episodeSimple := Episode{}

// 				episodeSimple.ID = episode.ID
// 				episodeSimple.EpisodeName = episode.EpisodeName
// 				episodeSimple.EpisodeNumber = episode.EpisodeNumber
// 				episodeSimple.FirstAired = episode.FirstAired
// 				episodeSimple.ImdbID = episode.ImdbID
// 				episodeSimple.Language = episode.Language
// 				episodeSimple.SeasonNumber = episode.SeasonNumber
// 				episodeSimple.LastUpdated = episode.LastUpdated
// 				episodeSimple.SeasonID = episode.SeasonID
// 				episodeSimple.ShowID = episode.ShowID

// 				if episode.FirstAired == "" {
// 					// fmt.Println("Missing first aired.")
// 				} else {
// 					aired, err := time.Parse("2006-01-02", episode.FirstAired)
// 					if err != nil {
// 						fmt.Println("Could not parse first aired.", err)
// 					} else {
// 						if aired.Before(time.Now()) {
// 							episodeSimple.HasAired = true
// 							// fmt.Println(series.SeriesName, "Season", episode.SeasonNumber, "Episode", episode.EpisodeNumber, "aired", episode.FirstAired)
// 						} else {
// 							episodeSimple.HasAired = false
// 							// fmt.Println(series.SeriesName, "Season", episode.SeasonNumber, "Episode", episode.EpisodeNumber, "not yet aired, airing on", episode.FirstAired)
// 						}
// 					}
// 				}
// 				// s.Episodes[SeasonEpisode{episode.SeasonNumber, episode.EpisodeNumber}] = &episodeSimple
// 				s.Episodes[fmt.Sprintf("%d_%d", episode.SeasonNumber, episode.EpisodeNumber)] = &episodeSimple
// 			}
// 		}
// 	} else {
// 		// TODO Log Error
// 	}
// }
