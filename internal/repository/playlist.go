package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type PlaylistStorage struct {
	db *sql.DB
}

func NewPlaylistStorage(db *sql.DB) (*PlaylistStorage, error) {
	return &PlaylistStorage{db: db}, nil
}

func (s *PlaylistStorage) Create(ctx context.Context, playlist *Playlist) error {
	const query = `INSERT INTO playlists (name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := s.db.ExecContext(
		ctx,
		query,
		playlist.Name,
		playlist.UserID,
		playlist.CreatedAt,
		playlist.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistStorage) Get(ctx context.Context, id int) (*Playlist, error) {
	const queryPlaylist = `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE id = $1`

	playlist := &Playlist{}
	err := s.db.QueryRowContext(ctx, queryPlaylist, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.UserID,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	const queryTracks = `
		SELECT t.id, t.name, t.artist, t.url, t.likes, t.dislikes
		FROM tracks t
		JOIN playlist_tracks pt ON pt.track_id = t.id
		WHERE pt.playlist_id = $1
	`

	trackRows, err := s.db.QueryContext(ctx, queryTracks, id)
	if err != nil {
		return nil, err
	}
	defer trackRows.Close()

	playlist.Tracks = []*Track{}

	for trackRows.Next() {
		track := &Track{}
		if err := trackRows.Scan(
			&track.ID,
			&track.Name,
			&track.Artist,
			&track.URL,
			&track.Likes,
			&track.Dislikes,
		); err != nil {
			return nil, err
		}
		playlist.Tracks = append(playlist.Tracks, track)
	}

	if err := trackRows.Err(); err != nil {
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistStorage) GetPlaylists(ctx context.Context, name string, userID int) ([]*Playlist, error) {
	baseQuery := `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE 1=1`
	var args []interface{}

	if name != "" {
		baseQuery += " AND name = $1"
		args = append(args, name)
	}
	if userID != 0 {
		baseQuery += " AND user_id = $2"
		args = append(args, userID)
	}

	rows, err := s.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []*Playlist
	for rows.Next() {
		playlist := &Playlist{}
		if err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.UserID,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

		const queryTracks = `
			SELECT t.id, t.name, t.artist, t.url, t.likes, t.dislikes
			FROM tracks t
			JOIN playlist_tracks pt ON pt.track_id = t.id
			WHERE pt.playlist_id = $1
		`
		trackRows, err := s.db.QueryContext(ctx, queryTracks, playlist.ID)
		if err != nil {
			return nil, err
		}
		defer trackRows.Close()

		var tracks []*Track
		for trackRows.Next() {
			track := &Track{}
			if err := trackRows.Scan(
				&track.ID,
				&track.Name,
				&track.Artist,
				&track.URL,
				&track.Likes,
				&track.Dislikes,
			); err != nil {
				return nil, err
			}
			tracks = append(tracks, track)
		}

		playlist.Tracks = tracks
		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func (s *PlaylistStorage) Update(ctx context.Context, playlist *Playlist) error {
	const query = `UPDATE playlists SET name = $1, user_id = $2, updated_at = $3 WHERE id = $4`

	_, err := s.db.ExecContext(ctx, query, playlist.Name, playlist.UserID, playlist.UpdatedAt, playlist.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM playlists WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistStorage) DeleteTracks(ctx context.Context, playlistID int, trackIDs []int) error {
	const query = `DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`

	if _, err := s.db.ExecContext(ctx, query, playlistID, pq.Array(trackIDs)); err != nil {
		return fmt.Errorf("failed to delete tracks: %w", err)
	}

	return nil
}

func (s *PlaylistStorage) AddTracks(ctx context.Context, playlistID int, trackIDs []int) error {
	const insertQuery = `INSERT INTO playlist_tracks (playlist_id, track_id) VALUES ($1, $2)`

	values := []string{}
	args := []interface{}{}
	for i, trackID := range trackIDs {
		values = append(values, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, trackID)
	}
	query := fmt.Sprintf(insertQuery, strings.Join(values, ", "))

	if _, err := s.db.ExecContext(ctx, query, append([]interface{}{playlistID}, args...)...); err != nil {
		return fmt.Errorf("failed to add tracks: %w", err)
	}

	return nil
}

func (s *PlaylistStorage) GetExistingTracks(ctx context.Context, playlistID int) ([]int, error) {
	const query = `SELECT track_id FROM playlist_tracks WHERE playlist_id = $1`
	rows, err := s.db.QueryContext(ctx, query, playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing tracks: %w", err)
	}
	defer rows.Close()

	var trackIDs []int
	for rows.Next() {
		var trackID int
		if err := rows.Scan(&trackID); err != nil {
			return nil, fmt.Errorf("failed to scan track ID: %w", err)
		}
		trackIDs = append(trackIDs, trackID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return trackIDs, nil
}
