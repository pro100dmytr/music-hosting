package services

import "time"

// TODO: move from /models/services/playlist.go to /models/playlist.go
type Playlist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	TracksID  []int     `json:"tracks_id"` // TODO: проверить что все работает
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
