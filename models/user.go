package models

type UserInterface interface {
}

type User struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Token    string `json:"token" db:"token"`
	Password string `json:"password" db:"password"`
}
