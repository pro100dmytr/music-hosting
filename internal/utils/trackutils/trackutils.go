package trackutils

import (
	"errors"
	"music-hosting/internal/models/services"
)

func ValidateTrack(track *services.Track) error {
	if track.Name == "" {
		return errors.New("name is required")
	}
	if track.Artist == "" {
		return errors.New("artist is required")
	}
	if track.URL == "" {
		return errors.New("url is required")
	}
	return nil
}
