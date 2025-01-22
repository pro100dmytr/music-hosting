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

// Todo: для create и update должны бьть разные модели. Создать CreatePlaylistRequest, UpdatelaylistRequest
type PlaylistRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// TODO: delete user id from request
	UserID   int   `json:"user_id"`
	TrackIDs []int `json:"tracks_id"`
}

type PlaylistResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	Tracks    []*Track  `json:"tracks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
