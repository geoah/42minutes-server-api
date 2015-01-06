package models

type Episode struct {
	ID             uint64
	EpisodeName    string
	EpisodeNumber  uint64
	FirstAired     string
	ImdbID         string
	Language       string
	SeasonNumber   uint64
	LastUpdated    string
	SeasonID       uint64
	SeriesID       uint64
	HasAired       bool
	LocalFilename  string
	LocalExists    bool
	LocalQuality   string
	TorrentQuality string
	TorrentLink    string
}
