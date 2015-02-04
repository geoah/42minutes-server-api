package models

import (
	"time"

	"github.com/42minutes/go-trakt"
)

type Episode struct {
	// ID                     int     `json:"id" db:"id"`
	ShowID                 int        `json:"show_id" db:"show_id"`
	TmdbID                 int        `json:"tmdb_id" db:"tmdb_id"`
	TraktID                int        `json:"trakt_id" db:"trakt_id"`
	TvdbID                 int        `json:"tvdb_id" db:"tvdb_id"`
	TvrageID               int        `json:"tvrage_id" db:"tvrage_id"`
	FirstAired             *time.Time `json:"first_aired" db:"first_aired"`
	ImagesScreenshotFull   string     `json:"images_screenshot_full" db:"images_screenshot_full"`
	ImagesScreenshotMedium string     `json:"images_screenshot_medium" db:"images_screenshot_medium"`
	ImagesScreenshotThumb  string     `json:"images_screenshot_thumb" db:"images_screenshot_thumb"`
	Episode                int        `json:"episode" db:"episode"`
	Overview               string     `json:"overview" db:"overview"`
	Rating                 float64    `json:"rating" db:"rating"`
	Season                 int        `json:"season" db:"season"`
	Title                  string     `json:"title" db:"title"`
	UpdatedAt              string     `json:"updated_at" db:"updated_at"`
	Votes                  int        `json:"votes" db:"votes"`
	Infohash               string     `json:"infohash" db:"infohash"`
	Available              bool       `json:"available" db:"-"`
	Watched                bool       `json:"watched" db:"-"`
}

func (e *Episode) MapInfo(traktEpisode trakt.Episode) {
	// e.ID = traktEpisode.IDs.Tvdb
	e.TmdbID = traktEpisode.IDs.Tmdb
	e.TraktID = traktEpisode.IDs.Trakt
	e.TvdbID = traktEpisode.IDs.Tvdb
	e.TvrageID = traktEpisode.IDs.Tvrage
	e.FirstAired = traktEpisode.FirstAired
	e.ImagesScreenshotFull = traktEpisode.Images.Screenshot.Full
	e.ImagesScreenshotMedium = traktEpisode.Images.Screenshot.Medium
	e.ImagesScreenshotThumb = traktEpisode.Images.Screenshot.Thumb
	e.Episode = traktEpisode.Number
	e.Overview = traktEpisode.Overview
	e.Rating = traktEpisode.Rating
	e.Season = traktEpisode.Season
	e.Title = traktEpisode.Title
	e.UpdatedAt = traktEpisode.UpdatedAt
	e.Votes = traktEpisode.Votes
}

// Personalize updates episode with additional information for a specific user
func (s *Episode) Normalize() {
	if s.FirstAired != nil {
		if time.Now().Sub(*s.FirstAired).Hours() > 0 {
			s.Status = "aired"
		} else {
			s.Status = "unaired"
		}
	}
}

// Personalize updates episode with additional information for a specific user
func (s *Episode) Personalize(userID string) {
	db := GetDbSession()
	userEpisode := UserEpisode{}
	err := db.SelectOne(&userEpisode, "select * from users_episodes where user_id=? and show_id=? and season=? and episode=", userID, s.ShowID, s.Season, s.Episode)
	if err == nil {
		s.Available = userEpisode.Available
		s.Watched = userEpisode.Watched
	}
}
