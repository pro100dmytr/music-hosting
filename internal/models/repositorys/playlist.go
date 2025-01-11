package repositorys

import "time"

// TODO: move to repository/models.go
type Playlist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	TracksID  []int     `json:"track_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
