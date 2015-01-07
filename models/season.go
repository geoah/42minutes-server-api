package models

type Season struct {
	ID       uint64    `json:"id" db:"id"`
	Season   int       `json:"season" db:"season"`
	URL      string    `json:"url" db:"url"`
	Poster   string    `json:"poster" db:"poster"`
	Episodes []Episode `json:"episodes" db:"-"`
}
