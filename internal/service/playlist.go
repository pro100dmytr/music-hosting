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

func (s *PlaylistService) CreatePlaylist(ctx context.Context, playlist *services.Playlist) error {
	if playlist.Name == "" {
		return fmt.Errorf("playlist name is required")
	}

	repoPlaylist := &repositorys.Playlist{
		Name:     playlist.Name,
		UserID:   playlist.UserID,
		TracksID: playlist.TracksID,
	}

	err := s.repo.Create(ctx, repoPlaylist)
	if err != nil {
		return err
	}

	if len(playlist.TracksID) > 0 {
		err = s.repo.AddTracksToPlaylist(ctx, playlist.ID, playlist.TracksID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PlaylistService) GetPlaylistByID(ctx context.Context, id int) (*services.Playlist, error) {
	repoPlaylist, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	tracksID, err := s.repo.GetTracksByPlaylistID(ctx, id)
	if err != nil {
		return nil, err
	}

	playlist := &services.Playlist{
		ID:        repoPlaylist.ID,
		Name:      repoPlaylist.Name,
		UserID:    repoPlaylist.UserID,
		TracksID:  tracksID,
		CreatedAt: repoPlaylist.CreatedAt,
		UpdatedAt: repoPlaylist.UpdatedAt,
	}

	return playlist, nil
}

func (s *PlaylistService) GetAllPlaylists(ctx context.Context) ([]*services.Playlist, error) {
	repoPlaylists, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var playlists []*services.Playlist
	for _, repoPlaylist := range repoPlaylists {
		tracksID, err := s.repo.GetTracksByPlaylistID(ctx, repoPlaylist.ID)
		if err != nil {
			return nil, err
		}

		playlist := &services.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			TracksID:  tracksID,
			CreatedAt: repoPlaylist.CreatedAt,
			UpdatedAt: repoPlaylist.UpdatedAt,
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (s *PlaylistService) GetPlaylistByName(ctx context.Context, name string) ([]*services.Playlist, error) {
	repoPlaylists, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var playlists []*services.Playlist
	for _, repoPlaylist := range repoPlaylists {
		tracksID, err := s.repo.GetTracksByPlaylistID(ctx, repoPlaylist.ID)
		if err != nil {
			return nil, err
		}

		playlist := &services.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			TracksID:  tracksID,
			CreatedAt: repoPlaylist.CreatedAt,
			UpdatedAt: repoPlaylist.UpdatedAt,
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (s *PlaylistService) GetPlaylistByUserID(ctx context.Context, userID int) ([]*services.Playlist, error) {
	repoPlaylists, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	var playlists []*services.Playlist
	for _, repoPlaylist := range repoPlaylists {
		tracksID, err := s.repo.GetTracksByPlaylistID(ctx, repoPlaylist.ID)
		if err != nil {
			return nil, err
		}

		playlist := &services.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			TracksID:  tracksID,
			CreatedAt: repoPlaylist.CreatedAt,
			UpdatedAt: repoPlaylist.UpdatedAt,
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, id int, playlist *services.Playlist) (*services.Playlist, error) {
	if playlist.Name == "" {
		return nil, fmt.Errorf("playlist name is required")
	}

	repoPlaylist := &repositorys.Playlist{
		Name:     playlist.Name,
		UserID:   playlist.UserID,
		TracksID: playlist.TracksID,
	}

	err := s.repo.Update(ctx, id, repoPlaylist)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdatePlaylistTracks(ctx, id, playlist.TracksID)
	if err != nil {
		return nil, err
	}

	playlist.ID = id
	return playlist, nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	err = s.repo.RemoveTracksFromPlaylist(ctx, id, []int{})
	if err != nil {
		return err
	}

	return nil
}
