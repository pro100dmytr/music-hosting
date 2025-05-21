package service

import (
	"context"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
)

type TrackService struct {
	trackRepo *repository.TrackStorage
	logger    *slog.Logger
}

func NewTrackService(trackRepo *repository.TrackStorage, logger *slog.Logger) *TrackService {
	return &TrackService{
		trackRepo: trackRepo,
		logger:    logger,
	}
}

func (s *TrackService) CreateTrack(ctx context.Context, track *models.Track) error {
	err := ValidateTrack(track)
	if err != nil {
		return err
	}

	repoTrack := repository.Track{
		Name:     track.Name,
		Artist:   track.Artist,
		URL:      track.URL,
		Likes:    track.Likes,
		Dislikes: track.Dislikes,
	}

	id, err := s.trackRepo.Create(ctx, &repoTrack)
	if err != nil {
		return err
	}

	track.ID = id
	return nil
}

func (s *TrackService) GetTrackByID(ctx context.Context, id int) (*models.Track, error) {
	repoTrack, err := s.trackRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	track := &models.Track{
		ID:       repoTrack.ID,
		Name:     repoTrack.Name,
		Artist:   repoTrack.Artist,
		URL:      repoTrack.URL,
		Likes:    repoTrack.Likes,
		Dislikes: repoTrack.Dislikes,
	}

	return track, nil
}

func (s *TrackService) UpdateTrack(ctx context.Context, track *models.Track) error {
	err := ValidateTrack(track)
	if err != nil {
		return err
	}

	trackRepo := repository.Track{
		ID:       track.ID,
		Name:     track.Name,
		Artist:   track.Artist,
		URL:      track.URL,
		Likes:    track.Likes,
		Dislikes: track.Dislikes,
	}

	err = s.trackRepo.Update(ctx, &trackRepo)
	if err != nil {
		return err
	}

	return nil
}

func (s *TrackService) DeleteTrack(ctx context.Context, id int) error {
	err := s.trackRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *TrackService) GetTracks(ctx context.Context, name, artist string, playlistID, offset, limit int) ([]*models.Track, error) {
	repoTracks, err := s.trackRepo.GetTracks(ctx, name, artist, playlistID, offset, limit)
	if err != nil {
		return nil, err
	}

	var tracks []*models.Track
	for _, repoTrack := range repoTracks {
		track := &models.Track{
			ID:       repoTrack.ID,
			Name:     repoTrack.Name,
			Artist:   repoTrack.Artist,
			URL:      repoTrack.URL,
			Likes:    repoTrack.Likes,
			Dislikes: repoTrack.Dislikes,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}
