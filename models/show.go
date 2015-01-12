package models

import (
	"fmt"
	"github.com/42minutes/go-trakt"
	"log"
)

type Show struct {
	ID                int      `json:"id" db:"id"`
	TmdbID            int      `json:"tmdb_id" db:"tmdb_id"`
	TraktID           int      `json:"trakt_id" db:"trakt_id"`
	TvdbID            int      `json:"tvdb_id" db:"tvdb_id"`
	TvrageID          int      `json:"tvrage_id" db:"tvrage_id"`
	AirDay            string   `json:"air_day" db:"air_day"`
	AirTime           string   `json:"air_time" db:"air_time"`
	AirTimezone       string   `json:"air_timezone" db:"air_timezone"`
	AiredEpisodes     int      `json:"aired_episodes" db:"aired_episodes"`
	Certification     string   `json:"certification" db:"certification"`
	Country           string   `json:"country" db:"country"`
	FirstAired        string   `json:"first_aired" db:"first_aired"`
	Homepage          string   `json:"homepage" db:"homepage"`
	ImageBannerFull   string   `json:"image_banner_full" db:"image_banner_full"`
	ImageClearartFull string   `json:"image_clearart_full" db:"image_clearart_full"`
	ImageFanartFull   string   `json:"image_fanart_full" db:"image_fanart_full"`
	ImageFanartMedium string   `json:"image_fanart_medium" db:"image_fanart_medium"`
	ImageFanartThumb  string   `json:"image_fanart_thumb" db:"image_fanart_thumb"`
	ImageLogoFull     string   `json:"image_logo_full" db:"image_logo_full"`
	ImagePosterFull   string   `json:"image_poster_full" db:"image_poster_full"`
	ImagePosterMedium string   `json:"image_poster_medium" db:"image_poster_medium"`
	ImagePosterThumb  string   `json:"image_poster_thumb" db:"image_poster_thumb"`
	ImageThumbFull    string   `json:"image_thumb_full" db:"image_thumb_full"`
	ImdbID            string   `json:"imdb_id" db:"imdb_id"`
	Language          string   `json:"language" db:"language"`
	Network           string   `json:"network" db:"network"`
	Overview          string   `json:"overview" db:"overview"`
	Rating            float64  `json:"rating" db:"rating"`
	Runtime           float64  `json:"runtime" db:"runtime"`
	Slug              string   `json:"slug" db:"slug"`
	Status            string   `json:"status" db:"status"`
	Title             string   `json:"title" db:"title"`
	Trailer           string   `json:"trailer" db:"trailer"`
	UpdatedAt         string   `json:"updated_at" db:"updated_at"`
	Votes             int      `json:"votes" db:"votes"`
	Year              int      `json:"year" db:"year"`
	Seasons           []Season `json:"seasons" db:"-"`
}

func (s *Show) MapInfo(traktShow trakt.Show) {
	s.ID = traktShow.Ids.Trakt
	s.TmdbID = traktShow.Ids.Tmdb
	s.TraktID = traktShow.Ids.Trakt
	s.TvdbID = traktShow.Ids.Tvdb
	s.TvrageID = traktShow.Ids.Tvrage
	s.AirDay = traktShow.Airs.Day
	s.AirTime = traktShow.Airs.Time
	s.AirTimezone = traktShow.Airs.Timezone
	s.AiredEpisodes = traktShow.AiredEpisodes
	s.Certification = traktShow.Certification
	s.Country = traktShow.Country
	s.FirstAired = traktShow.FirstAired
	s.Homepage = traktShow.Homepage
	s.ImageBannerFull = traktShow.Images.Banner.Full
	s.ImageClearartFull = traktShow.Images.Clearart.Full
	s.ImageFanartFull = traktShow.Images.Fanart.Full
	s.ImageFanartMedium = traktShow.Images.Fanart.Medium
	s.ImageFanartThumb = traktShow.Images.Logo.Full
	s.ImageLogoFull = traktShow.Images.Poster.Full
	s.ImagePosterFull = traktShow.Images.Poster.Full
	s.ImagePosterMedium = traktShow.Images.Poster.Medium
	s.ImagePosterThumb = traktShow.Images.Poster.Thumb
	s.ImageThumbFull = traktShow.Images.Thumb.Full
	s.ImdbID = traktShow.Ids.Imdb
	s.Language = traktShow.Language
	s.Network = traktShow.Network
	s.Overview = traktShow.Overview
	s.Rating = traktShow.Rating
	s.Runtime = traktShow.Runtime
	s.Slug = traktShow.Ids.Slug
	s.Status = traktShow.Status
	s.Title = traktShow.Title
	s.Trailer = traktShow.Trailer
	s.UpdatedAt = traktShow.UpdatedAt
	s.Votes = traktShow.Votes
	s.Year = traktShow.Year

	s.Seasons = make([]Season, 0)
	// for _, aseason := range traktShow.Seasons {
	// 	season := Season{}
	// 	season.ShowID = s.ID
	// 	season.MapInfo(aseason)
	// 	s.Seasons = append(s.Seasons, season)
	// }
}

// Get basic information from tvdbcom by TvdbID
func (s *Show) UpdateInfoByTraktID(traktID int) {
	var trakt trakt.Client = *GetTraktSession()
	log.Printf("Trying to retrieve information for show:traktid:%d", traktID)
	show, result := trakt.Shows().One(traktID)
	fmt.Println("SHOWXX", show)
	if result.HasError() == false {
		s.MapInfo(*show)
	} else {
		fmt.Println("ERROR", result.Error())
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
