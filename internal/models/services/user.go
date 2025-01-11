package services

// TODO: move from /models/services/user.go to /models/user.go
type User struct {
	ID         int
	Login      string
	Email      string
	Password   string
	PlaylistID []int // TODO: проверить что все работает
}
