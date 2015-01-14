package models

type ShowMatch struct {
	Title  string `json:"show_id" db:"title"`
	ShowID int    `json:"show_id" db:"show_id"`
}
