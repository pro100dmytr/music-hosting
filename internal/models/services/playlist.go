package services

import "time"

type Playlist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	TracksID  []int     `json:"tracks_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
