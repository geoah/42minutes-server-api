package models

import "github.com/42minutes/go-trakt"

type Season struct {
	// ID                 int       `json:"id" db:"id"`
	ShowID             int       `json:"show_id" db:"show_id"`
	TmdbID             int       `json:"tmdb_id" db:"tmdb_id"`
	TraktID            int       `json:"trakt_id" db:"trakt_id"`
	TvdbID             int       `json:"tvdb_id" db:"tvdb_id"`
	TvrageID           int       `json:"tvrage_id" db:"tvrage_id"`
	EpisodeCount       int       `json:"episode_count" db:"episode_count"`
	ImagesPosterFull   string    `json:"images_poster_full" db:"images_poster_full"`
	ImagesPosterMedium string    `json:"images_poster_medium" db:"images_poster_medium"`
	ImagesPosterThumb  string    `json:"images_poster_thumb" db:"images_poster_thumb"`
	ImagesThumbFull    string    `json:"images_thumb_full" db:"images_thumb_full"`
	Season             int       `json:"season" db:"season"`
	Overview           string    `json:"overview" db:"overview"`
	Rating             float64   `json:"rating" db:"rating"`
	Votes              int       `json:"votes" db:"votes"`
	Episodes           []Episode `json:"episodes" db:"-"`
}

func (s *Season) MapInfo(traktSeason trakt.Season) {
	// s.ID = traktSeason.Ids.Tvdb
	s.TmdbID = traktSeason.IDs.Tmdb
	s.TraktID = traktSeason.IDs.Trakt
	s.TvdbID = traktSeason.IDs.Tvdb
	s.TvrageID = traktSeason.IDs.Tvrage
	s.EpisodeCount = traktSeason.EpisodeCount
	s.ImagesPosterFull = traktSeason.Images.Poster.Full
	s.ImagesPosterMedium = traktSeason.Images.Poster.Medium
	s.ImagesPosterThumb = traktSeason.Images.Poster.Thumb
	s.ImagesThumbFull = traktSeason.Images.Thumb.Full
	s.Season = traktSeason.Number
	s.Overview = traktSeason.Overview
	s.Rating = traktSeason.Rating
	s.Votes = traktSeason.Votes
}
