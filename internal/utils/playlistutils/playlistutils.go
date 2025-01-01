package playlistutils

import (
	"fmt"
	"music-hosting/internal/models/services"
)

func ValidatePlaylist(playlist *services.Playlist) error {
	if playlist.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}
