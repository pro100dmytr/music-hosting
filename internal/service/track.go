package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
	"strconv"
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

func ValidateTrack(track *models.Track) error {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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

func (s *TrackService) GetAllTracks(ctx context.Context) ([]*models.Track, error) {
	repoTracks, err := s.trackRepo.GetAll(ctx)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return nil
}

func (s *TrackService) DeleteTrack(ctx context.Context, id int) error {
	err := s.trackRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	return nil
}

func (s *TrackService) GetTracksWithPagination(ctx context.Context, offset, limit string) ([]*models.Track, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		return nil, fmt.Errorf("invalid limit parametr")
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		return nil, err
	}

	repoTracks, err := s.trackRepo.GetTracks(ctx, offsetInt, limitInt)
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

func (s *TrackService) GetTracksByName(ctx context.Context, name string) ([]*models.Track, error) {
	repoTracks, err := s.trackRepo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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

func (s *TrackService) GetTracksByArtist(ctx context.Context, artist string) ([]*models.Track, error) {
	repoTracks, err := s.trackRepo.GetByArtist(ctx, artist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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

func (s *TrackService) GetTracksByPlaylistID(ctx context.Context, playlistID int) ([]*models.Track, error) {
	repoTracks, err := s.trackRepo.GetTracksByPlaylistID(ctx, playlistID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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
