package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"music-hosting/internal/models/repositorys"
	"music-hosting/internal/models/services"
	"music-hosting/internal/repository"
	"music-hosting/internal/utils/trackutils"
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

func (s *TrackService) CreateTrack(ctx context.Context, track *services.Track) error {
	err := trackutils.ValidateTrack(track)
	if err != nil {
		return err
	}

	repoTrack := repositorys.Track{
		Name:   track.Name,
		Artist: track.Artist,
		URL:    track.URL,
	}

	id, err := s.trackRepo.Create(ctx, &repoTrack)
	if err != nil {
		return err
	}

	track.ID = id
	return nil
}

func (s *TrackService) GetTrackByID(ctx context.Context, id int) (*services.Track, error) {
	repoTrack, err := s.trackRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	track := &services.Track{
		ID:     repoTrack.ID,
		Name:   repoTrack.Name,
		Artist: repoTrack.Artist,
		URL:    repoTrack.URL,
	}

	return track, nil
}

func (s *TrackService) GetAllTracks(ctx context.Context) ([]*services.Track, error) {
	repoTracks, err := s.trackRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var tracks []*services.Track
	for _, repoTrack := range repoTracks {
		track := &services.Track{
			ID:     repoTrack.ID,
			Name:   repoTrack.Name,
			Artist: repoTrack.Artist,
			URL:    repoTrack.URL,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (s *TrackService) UpdateTrack(ctx context.Context, id int, track *services.Track) (*services.Track, error) {
	err := trackutils.ValidateTrack(track)
	if err != nil {
		return nil, err
	}

	trackRepo := repositorys.Track{
		ID:     track.ID,
		Name:   track.Name,
		Artist: track.Artist,
		URL:    track.URL,
	}

	err = s.trackRepo.Update(ctx, &trackRepo, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	track.ID = id
	return track, nil
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

func (s *TrackService) GetTracksWithPagination(ctx context.Context, offset, limit string) ([]*services.Track, error) {
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

	var tracks []*services.Track
	for _, repoTrack := range repoTracks {
		track := &services.Track{
			ID:     repoTrack.ID,
			Name:   repoTrack.Name,
			Artist: repoTrack.Artist,
			URL:    repoTrack.URL,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (s *TrackService) GetTrackByName(ctx context.Context, name string) ([]*services.Track, error) {
	repoTracks, err := s.trackRepo.GetForName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var tracks []*services.Track
	for _, repoTrack := range repoTracks {
		track := &services.Track{
			ID:     repoTrack.ID,
			Name:   repoTrack.Name,
			Artist: repoTrack.Artist,
			URL:    repoTrack.URL,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (s *TrackService) GetTrackByArtist(ctx context.Context, artist string) ([]*services.Track, error) {
	repoTracks, err := s.trackRepo.GetForArtist(ctx, artist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var tracks []*services.Track
	for _, repoTrack := range repoTracks {
		track := &services.Track{
			ID:     repoTrack.ID,
			Name:   repoTrack.Name,
			Artist: repoTrack.Artist,
			URL:    repoTrack.URL,
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}
