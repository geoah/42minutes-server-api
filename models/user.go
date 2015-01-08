package models

type User struct {
	ID    string `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	Token int    `json:"token" db:"token"`
}