package models

type User struct {
	ID       int        `json:"id"`
	Login    string     `json:"login"`
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Playlist []Playlist `json:"-"`
}