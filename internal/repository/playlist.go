package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (s *PlaylistStorage) GetAll(ctx context.Context) ([]*Playlist, error) {
	// TODO: удалить ненужную транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	const queryPlaylists = `SELECT id, name, user_id, created_at, updated_at FROM playlists`
	rows, err := s.db.QueryContext(ctx, queryPlaylists)
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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return playlists, nil
}

func (s *PlaylistStorage) Get(ctx context.Context, id int) (*Playlist, error) {
	// TODO: удалить ненужную транзакцию

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	const queryPlaylist = `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE id = $1`

	playlist := &Playlist{}
	err = s.db.QueryRowContext(ctx, queryPlaylist, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.UserID,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Плейлист не найден
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

	if err := trackRows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return playlist, nil
}
func (s *PlaylistStorage) GetByUserID(ctx context.Context, userID int) ([]*Playlist, error) {
	// TODO: удалить ненужную транзакцию

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	const queryPlaylists = `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, queryPlaylists, userID)
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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return playlists, nil
}

func (s *PlaylistStorage) GetByName(ctx context.Context, name string) ([]*Playlist, error) {
	// TODO: удалить ненужную транзакцию

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	const queryPlaylists = `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, queryPlaylists, name)
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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return playlists, nil
}

func (s *PlaylistStorage) Update(ctx context.Context, playlist *Playlist) error {
	const query = `UPDATE playlists SET name = $1, user_id = $2, updated_at = $3 WHERE id = $4`

	result, err := s.db.ExecContext(ctx, query, playlist.Name, playlist.UserID, playlist.UpdatedAt, playlist.ID)
	if err != nil {
		return err
	}

	// TODO: если не проапдейтили то это не ошибка
	//rowsAffected, err := result.RowsAffected()
	//if err != nil {
	//	return err
	//}
	//
	//if rowsAffected == 0 {
	//	return sql.ErrNoRows
	//}

	return nil
}

func (s *PlaylistStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM playlists WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// TODO: DELETE
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

// TODO: Разделить UpdatePlaylistTracks на AddTracks и DeleteTracks
// а логику определения того что нужно добавить или удалить нужно вынести
// в сервис
func (s *PlaylistStorage) UpdatePlaylistTracks(ctx context.Context, playlistID int, trackIDs []int) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	existingTrackIDs, err := s.getExistingTracks(ctx, tx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to get existing tracks: %w", err)
	}

	toDelete := difference(existingTrackIDs, trackIDs)
	toAdd := difference(trackIDs, existingTrackIDs)

	if len(toDelete) > 0 {
		const deleteQuery = `DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`
		for _, trackID := range toDelete {
			if _, err := tx.ExecContext(ctx, deleteQuery, playlistID, trackID); err != nil {
				return fmt.Errorf("failed to delete track %d: %w", trackID, err)
			}
		}
	}

	if len(toAdd) > 0 {
		const insertQuery = `INSERT INTO playlist_tracks (playlist_id, track_id) VALUES ($1, $2)`
		for _, trackID := range toAdd {
			if _, err := tx.ExecContext(ctx, insertQuery, playlistID, trackID); err != nil {
				return fmt.Errorf("failed to add track %d: %w", trackID, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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

func (s *PlaylistStorage) getExistingTracks(ctx context.Context, tx *sql.Tx, playlistID int) ([]int, error) {
	const selectQuery = `SELECT track_id FROM playlist_tracks WHERE playlist_id = $1`
	rows, err := tx.QueryContext(ctx, selectQuery, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var existingTrackIDs []int
	for rows.Next() {
		var trackID int
		if err := rows.Scan(&trackID); err != nil {
			return nil, err
		}
		existingTrackIDs = append(existingTrackIDs, trackID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return existingTrackIDs, nil
}

func (s *PlaylistStorage) GetPlaylistTrackCount(ctx context.Context, playlistID int) (int, error) {
	const countQuery = `SELECT COUNT(*) FROM playlist_tracks WHERE playlist_id = $1`
	var count int
	err := s.db.QueryRowContext(ctx, countQuery, playlistID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get track count: %w", err)
	}
	return count, nil
}
