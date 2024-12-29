package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
	"music-hosting/pkg/utils/trackutils"
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

func (s *TrackService) CreateTrack(ctx context.Context, track *models.Track) (*models.Track, error) {
	if err := trackutils.ValidateTrack(track); err != nil {
		s.logger.Error("Error validating track", slog.Any("error", err))
		return nil, err
	}

	if err := s.trackRepo.Create(ctx, track); err != nil {
		s.logger.Error("Error creating track in database", slog.Any("error", err))
		return nil, err
	}

	return track, nil
}

func (s *TrackService) GetTrackByID(ctx context.Context, id int) (*models.Track, error) {
	track, err := s.trackRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("Track not found", slog.Any("error", err))
			return nil, err
		}
		s.logger.Error("Error fetching track by ID", slog.Any("error", err))
		return nil, err
	}
	return track, nil
}

func (s *TrackService) GetAllTracks(ctx context.Context) ([]*models.Track, error) {
	tracks, err := s.trackRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Error fetching all tracks", slog.Any("error", err))
		return nil, err
	}
	return tracks, nil
}

func (s *TrackService) GetTrackByName(ctx context.Context, name string) ([]*models.Track, error) {
	track, err := s.trackRepo.GetForName(ctx, name)
	if err != nil {
		s.logger.Error("Error fetching track by name", slog.Any("error", err))
		return nil, err
	}
	return track, nil
}

func (s *TrackService) GetTrackByArtist(ctx context.Context, artist string) ([]*models.Track, error) {
	track, err := s.trackRepo.GetForArtist(ctx, artist)
	if err != nil {
		s.logger.Error("Error fetching track by artist", slog.Any("error", err))
		return nil, err
	}
	return track, nil
}

func (s *TrackService) UpdateTrack(ctx context.Context, track *models.Track, id int) error {
	if err := trackutils.ValidateTrack(track); err != nil {
		s.logger.Error("Error validating track", slog.Any("error", err))
		return err
	}

	if err := s.trackRepo.Update(ctx, track, id); err != nil {
		s.logger.Error("Error updating track in database", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *TrackService) DeleteTrack(ctx context.Context, id int) error {
	if err := s.trackRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("Track not found for deletion", slog.Any("error", err))
			return err
		}
		s.logger.Error("Error deleting track", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *TrackService) GetTracksWithPagination(ctx context.Context, limit, offset int) ([]*models.Track, error) {
	tracks, err := s.trackRepo.GetTracks(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Error fetching tracks with pagination", slog.Any("error", err))
		return nil, err
	}
	return tracks, nil
}
