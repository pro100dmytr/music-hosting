package models

import "time"

type Playlist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	Tracks    []*Track  `json:"tracks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlaylistRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserID   int    `json:"user_id"`
	TrackIDs []int  `json:"tracks_id"`
}

type PlaylistResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	Tracks    []*Track  `json:"tracks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
