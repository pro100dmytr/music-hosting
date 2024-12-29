package service

import (
	"context"
	"errors"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
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

func (s *PlaylistService) CreatePlaylist(ctx context.Context, playlist *models.Playlist) (*models.Playlist, error) {
	if playlist.Name == "" {
		s.logger.Error("Playlist name is required")
		return nil, errors.New("playlist name is required")
	}

	err := s.repo.Create(ctx, playlist)
	if err != nil {
		s.logger.Error("Failed to create playlist", slog.Any("error", err))
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistService) GetPlaylistByID(ctx context.Context, id int) (*models.Playlist, error) {
	playlist, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get playlist", slog.Any("error", err))
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistService) GetAllPlaylists(ctx context.Context) ([]*models.Playlist, error) {
	playlists, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get playlists", slog.Any("error", err))
		return nil, err
	}

	return playlists, nil
}

func (s *PlaylistService) GetPlaylistByName(ctx context.Context, name string) ([]*models.Playlist, error) {
	playlists, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error("Failed to get playlists", slog.Any("error", err))
		return nil, err
	}

	return playlists, nil
}

func (s *PlaylistService) GetPlaylistByUserID(ctx context.Context, userID int) ([]*models.Playlist, error) {
	playlists, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get playlists", slog.Any("error", err))
		return nil, err
	}

	return playlists, nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, id int, playlist *models.Playlist) (*models.Playlist, error) {
	err := s.repo.Update(ctx, id, playlist)
	if err != nil {
		s.logger.Error("Failed to update playlist", slog.Any("error", err))
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete playlist", slog.Any("error", err))
		return err
	}

	return nil
}
