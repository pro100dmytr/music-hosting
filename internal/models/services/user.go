package services

type User struct {
	ID         int
	Login      string
	Email      string
	Password   string
	PlaylistID []int
}
