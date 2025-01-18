package repository

import (
	"music-hosting/internal/models"
	"time"
)

type User struct {
	ID       int
	Login    string
	Email    string
	Password string
	Salt     string
}

type Track struct {
	ID       int
	Name     string
	Artist   string
	URL      string
	Likes    int
	Dislikes int
}

// TODO: delete JSON tags
type Playlist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	Tracks    []*Track  `json:"tracks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Track) ConvertToModel() *models.Track {
	return &models.Track{
		ID:       t.ID,
		Name:     t.Name,
		Artist:   t.Artist,
		URL:      t.URL,
		Likes:    t.Likes,
		Dislikes: t.Dislikes,
	}
}
