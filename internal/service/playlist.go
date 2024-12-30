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
	"music-hosting/pkg/utils/convertutils"
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

	repoPlaylist := repositorys.Playlist{
		Name:    playlist.Name,
		UserID:  playlist.UserID,
		TrackID: convertutils.IntSliceConvertIntoString(playlist.TrackID),
	}

	err := s.repo.Create(ctx, &repoPlaylist)
	if err != nil {
		return err
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

	trackID := []int{}
	if repoPlaylist.TrackID != "" {
		trackID, err = convertutils.StringConvertIntoIntSlice(repoPlaylist.TrackID)
		if err != nil {
			return nil, err
		}
	}

	playlist := &services.Playlist{
		ID:        repoPlaylist.ID,
		Name:      repoPlaylist.Name,
		UserID:    repoPlaylist.UserID,
		TrackID:   trackID,
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
		trackID := []int{}
		if repoPlaylist.TrackID != "" {
			trackID, err = convertutils.StringConvertIntoIntSlice(repoPlaylist.TrackID)
			if err != nil {
				return nil, err
			}
		}

		playlist := &services.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			TrackID:   trackID,
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
		trackID := []int{}
		if repoPlaylist.TrackID != "" {
			trackID, err = convertutils.StringConvertIntoIntSlice(repoPlaylist.TrackID)
			if err != nil {
				return nil, err
			}
		}

		playlist := &services.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			TrackID:   trackID,
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
		trackID := []int{}
		if repoPlaylist.TrackID != "" {
			trackID, err = convertutils.StringConvertIntoIntSlice(repoPlaylist.TrackID)
			if err != nil {
				return nil, err
			}
		}

		playlist := &services.Playlist{
			ID:        repoPlaylist.ID,
			Name:      repoPlaylist.Name,
			UserID:    repoPlaylist.UserID,
			TrackID:   trackID,
			CreatedAt: repoPlaylist.CreatedAt,
			UpdatedAt: repoPlaylist.UpdatedAt,
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, id int, playlist *services.Playlist) (*services.Playlist, error) {
	stringIDString := convertutils.IntSliceConvertIntoString(playlist.TrackID)
	repoPlaylist := &repositorys.Playlist{
		Name:    playlist.Name,
		UserID:  playlist.UserID,
		TrackID: stringIDString,
	}

	err := s.repo.Update(ctx, id, repoPlaylist)
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

	return nil
}
