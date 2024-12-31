package playlistutils

import (
	"errors"
	"fmt"
	"music-hosting/internal/models/services"
)

func ValidatePlaylist(playlist *services.Playlist) error {
	if playlist.Name == "" {
		return fmt.Errorf("name is required")
	}
	if track.Artist == "" {
		return errors.New("artist is required")
	}
	if track.URL == "" {
		return errors.New("url is required")
	}
	return nil
}

