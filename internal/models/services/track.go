package services

// TODO: move from /models/services/track.go to /models/track.go
type Track struct {
	ID       int
	Name     string
	Artist   string
	URL      string
	Likes    int
	Dislikes int
}
