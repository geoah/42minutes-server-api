package models

import (
	"fmt"
	"github.com/hobeone/gotrakt"
	"log"
	// "time"
)

type Show struct {
	ID            int               `json:"id" db:"id"`
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
	TvdbID        int               `json:"tvdb_id" db:"tvdb_id"`
	TvrageID      int               `json:"tvrage_id" db:"tvrage_id"`
	Ended         bool              `json:"ended" db:"ended"`
	Images        map[string]string `json:"images" db:"-"`
	Genres        []string          `json:"genres" db:"-"`
	Seasons       []Season          `json:"seasons" db:"-"`
}

func (s *Show) mapInfo(show gotrakt.Show) {
	s.ID = show.TvdbID
	s.Title = show.Title
	s.Year = show.Year
	s.URL = show.URL
	s.FirstAired = show.FirstAired
	s.Country = show.Country
	s.Overview = show.Overview
	s.Runtime = show.Runtime
	s.Network = show.Network
	s.AirDay = show.AirDay
	s.AirTime = show.AirTime
	s.Certification = show.Certification
	s.ImdbID = show.ImdbID
	s.TvdbID = show.TvdbID
	s.TvrageID = show.TvrageID
	s.Ended = show.Ended
	s.Images = show.Images
	s.Genres = show.Genres

	// s.Seasons = make([]Season)
	// for season := range show.Seasons {

	// }
}

// Get basic information from tvdbcom by TvdbID
func (s *Show) UpdateInfoByTvdbID(tvdbID int) {
	var trakt gotrakt.TraktTV = *GetTraktSession()
	log.Printf("Trying to retrieve information for show:tvdbid:%d", tvdbID)
	show, err := trakt.GetShow(fmt.Sprintf("%d", tvdbID))
	if err == nil {
		s.mapInfo(*show)
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
