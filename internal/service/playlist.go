package service

import (
	"context"
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
		return nil, err
	}

	if repoPlaylist == nil {
		return nil, nil
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
func (s *PlaylistService) GetPlaylists(ctx context.Context, name string, userID int) ([]*models.Playlist, error) {
	repoPlaylists, err := s.repo.GetPlaylists(ctx, name, userID)
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

func (s *PlaylistService) UpdatePlaylistTracks(ctx context.Context, playlistID int, newTrackIDs []int) error {
	existingTrackIDs, err := s.repo.GetExistingTracks(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to get existing tracks: %w", err)
	}

	toDelete := difference(existingTrackIDs, newTrackIDs)
	toAdd := difference(newTrackIDs, existingTrackIDs)

	if len(toDelete) > 0 {
		if err := s.repo.DeleteTracks(ctx, playlistID, toDelete); err != nil {
			return fmt.Errorf("failed to delete tracks: %w", err)
		}
	}

	if len(toAdd) > 0 {
		if err := s.repo.AddTracks(ctx, playlistID, toAdd); err != nil {
			return fmt.Errorf("failed to add tracks: %w", err)
		}
	}

	return nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, playlist *models.Playlist, trackIDs []int) error {
	if playlist.Name == "" {
		return fmt.Errorf("playlist name is required")
	}

	repoPlaylist := &repository.Playlist{
		ID:        playlist.ID,
		Name:      playlist.Name,
		UserID:    playlist.UserID,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.Update(ctx, repoPlaylist); err != nil {
		return fmt.Errorf("failed to update playlist: %w", err)
	}

	if err := s.UpdatePlaylistTracks(ctx, repoPlaylist.ID, trackIDs); err != nil {
		return fmt.Errorf("failed to update playlist tracks: %w", err)
	}

	return nil
}

func difference(slice1, slice2 []int) []int {
	m := make(map[int]struct{})
	for _, v := range slice2 {
		m[v] = struct{}{}
	}

	var diff []int
	for _, v := range slice1 {
		if _, found := m[v]; !found {
			diff = append(diff, v)
		}
	}
	return diff
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
