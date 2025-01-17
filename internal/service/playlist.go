package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
	"time"
)

type PlaylistService struct {
	repo   *repository.PlaylistStorage
	logger *slog.Logger
}

func NewPlaylistService(repo *repository.PlaylistStorage, logger *slog.Logger) *PlaylistService {
	return &PlaylistService{
		repo:   repo,
		logger: logger,
	}
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	if playlist.Name == "" {
		return fmt.Errorf("playlist name is required")
	}

	repoPlaylist := &repository.Playlist{
		Name:      playlist.Name,
		UserID:    playlist.UserID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := s.repo.Create(ctx, repoPlaylist)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistService) GetPlaylistByID(ctx context.Context, id int) (*models.Playlist, error) {
	repoPlaylist, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var tracks []*models.Track
	for _, track := range repoPlaylist.Tracks {
		tracks = append(tracks, track.ConvertToModel())
	}

	playlist := &models.Playlist{
		ID:        repoPlaylist.ID,
		Name:      repoPlaylist.Name,
		UserID:    repoPlaylist.UserID,
		Tracks:    tracks,
		CreatedAt: repoPlaylist.CreatedAt,
		UpdatedAt: repoPlaylist.UpdatedAt,
	}

	return playlist, nil
}

func (s *PlaylistService) GetAllPlaylists(ctx context.Context) ([]*models.Playlist, error) {
	repoPlaylists, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var playlists []*models.Playlist
	for _, repoPlaylist := range repoPlaylists {
		playlist := &models.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			CreatedAt: repoPlaylist.CreatedAt,
			UpdatedAt: repoPlaylist.UpdatedAt,
		}

		var tracks []*models.Track
		for _, repoTrack := range repoPlaylist.Tracks {
			tracks = append(tracks, repoTrack.ConvertToModel())
		}
		playlist.Tracks = tracks

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (s *PlaylistService) GetPlaylistsByName(ctx context.Context, name string) ([]*models.Playlist, error) {
	repoPlaylists, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var playlists []*models.Playlist
	for _, repoPlaylist := range repoPlaylists {
		playlist := &models.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			CreatedAt: repoPlaylist.CreatedAt,
			UpdatedAt: repoPlaylist.UpdatedAt,
		}

		var tracks []*models.Track
		for _, repoTrack := range repoPlaylist.Tracks {
			tracks = append(tracks, repoTrack.ConvertToModel())
		}
		playlist.Tracks = tracks

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (s *PlaylistService) GetPlaylistsByUserID(ctx context.Context, userID int) ([]*models.Playlist, error) {
	repoPlaylists, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var playlists []*models.Playlist
	for _, repoPlaylist := range repoPlaylists {
		playlist := &models.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			CreatedAt: repoPlaylist.CreatedAt,
			UpdatedAt: repoPlaylist.UpdatedAt,
		}

		var tracks []*models.Track
		for _, repoTrack := range repoPlaylist.Tracks {
			tracks = append(tracks, repoTrack.ConvertToModel())
		}
		playlist.Tracks = tracks

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, playlist *models.Playlist, trackIDs []int) error {
	if playlist.Name == "" {
		return fmt.Errorf("playlist name is required")
	}

	repoPlaylist := &repository.Playlist{
		Name:      playlist.Name,
		UserID:    playlist.UserID,
		UpdatedAt: time.Now().UTC(),
	}

	err := s.repo.Update(ctx, repoPlaylist)
	if err != nil {
		return err
	}

	existingTrackCount, err := s.repo.GetPlaylistTrackCount(ctx, repoPlaylist.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing track count: %w", err)
	}

	if existingTrackCount == len(trackIDs) {
		return fmt.Errorf("no changes detected, playlist is already up to date")
	}

	err = s.repo.UpdatePlaylistTracks(ctx, repoPlaylist.ID, trackIDs)
	if err != nil {
		return fmt.Errorf("failed to update playlist tracks: %w", err)
	}

	return nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	return nil
}
