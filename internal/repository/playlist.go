package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music-hosting/internal/config"
	"music-hosting/internal/database/postgresql"
	"music-hosting/internal/models/repositorys"

	"github.com/lib/pq"
)

type PlaylistStorage struct {
	db *sql.DB
}

func (s *PlaylistStorage) Close() error {
	return postgresql.CloseConn(s.db)
}

func NewPlaylistStorage(cfg *config.Config) (*PlaylistStorage, error) {
	db, err := postgresql.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &PlaylistStorage{db: db}, nil
}

func (s *PlaylistStorage) Create(ctx context.Context, playlist *repositorys.Playlist) error {
	const query = `INSERT INTO playlists (name, user_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	// TODO: remove .Scan
	err := s.db.QueryRowContext(
		ctx,
		query,
		playlist.Name,
		playlist.UserID,
	).Scan(
		&playlist.ID,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistStorage) GetAll(ctx context.Context) ([]*repositorys.Playlist, error) {
	const query = `SELECT id, name, user_id, created_at, updated_at FROM playlists`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []*repositorys.Playlist
	for rows.Next() {
		playlist := &repositorys.Playlist{}
		if err := rows.Scan(
			playlist.ID,
			playlist.Name,
			playlist.UserID,
			playlist.CreatedAt,
			playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func (s *PlaylistStorage) Get(ctx context.Context, id int) (*repositorys.Playlist, error) {
	const query = `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE id = $1`

	playlist := &repositorys.Playlist{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		playlist.ID,
		playlist.Name,
		playlist.UserID,
		playlist.CreatedAt,
		playlist.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows // TODO: Return nil instead of sql.ErrNoRows
		}
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistStorage) GetByUserID(ctx context.Context, userID int) ([]*repositorys.Playlist, error) {
	const query = `SELECT * FROM playlists WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var playlists []*repositorys.Playlist
	for rows.Next() {
		playlist := &repositorys.Playlist{}
		if err := rows.Scan(
			playlist.ID,
			playlist.Name,
			playlist.UserID,
			playlist.CreatedAt,
			playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, rows.Err()
}

func (s *PlaylistStorage) GetByName(ctx context.Context, name string) ([]*repositorys.Playlist, error) {
	const query = `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var playlists []*repositorys.Playlist
	for rows.Next() {
		playlist := &repositorys.Playlist{}
		if err := rows.Scan(
			playlist.ID,
			playlist.Name,
			playlist.UserID,
			playlist.CreatedAt,
			playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, rows.Err()
}

func (s *PlaylistStorage) Update(ctx context.Context, id int, playlist *repositorys.Playlist) error {
	// TODO: тут не нужна транзакция. Удали ее
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `UPDATE playlists SET name = $1, user_id = $2, updated_at = NOW() WHERE id = $3`

	result, err := s.db.ExecContext(ctx, query, playlist.Name, playlist.UserID, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PlaylistStorage) Delete(ctx context.Context, id int) error {
	// TODO: тут не нужна транзакция. Удали ее
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `DELETE FROM playlists WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// TODO: delete statement
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TODO: Лучше удалить этот метод и в TrackRepository создать метод который возвращает
// список треков по айди плейлиста
// func () GetTracksByPlaylistID(...) ([]Track, err)
func (s *PlaylistStorage) GetTracksByPlaylistID(ctx context.Context, playlistID int) ([]int, error) {
	const query = `SELECT track_id FROM playlist_tracks WHERE playlist_id = $1`

	rows, err := s.db.QueryContext(ctx, query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracksID []int
	for rows.Next() {
		var trackID int
		if err := rows.Scan(&trackID); err != nil {
			return nil, err
		}
		tracksID = append(tracksID, trackID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracksID, nil
}

// TODO: функционал реализован неверно. Не нцжно сначала удалять все а потом вставлять.
// нужно на уровне сервиса определить каких треков нету в плейлисте и добавих их или удалить
func (s *PlaylistStorage) UpdatePlaylistTracks(ctx context.Context, playlistID int, trackIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const deleteQuery = `DELETE FROM playlist_tracks WHERE playlist_id = $1`
	_, err = s.db.ExecContext(ctx, deleteQuery, playlistID)
	if err != nil {
		return err
	}

	const insertQuery = `INSERT INTO playlist_tracks (playlist_id, track_id) VALUES ($1, $2)`
	for _, trackID := range trackIDs {
		_, err := s.db.ExecContext(ctx, insertQuery, playlistID, trackID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TODO: функция не работает. Нужно переписать
func (s *PlaylistStorage) AddTracksToPlaylist(ctx context.Context, playlistID int, trackIDs []int) error {
	// TODO: удалить транзакцию если не нужна
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	if playlistID <= 0 || len(trackIDs) == 0 {
		return fmt.Errorf("invalid input: playlistID or trackIDs are empty")
	}

	const query = `
		INSERT INTO playlist_tracks (playlist_id, track_id)
		VALUES ($1, unnest($2::int[]))
	`
	_, err = s.db.ExecContext(ctx, query, playlistID, pq.Array(trackIDs))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TODO: функция не работает. Нужно переписать
func (s *PlaylistStorage) RemoveTracksFromPlaylist(ctx context.Context, playlistID int, trackIDs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	if playlistID <= 0 {
		return fmt.Errorf("invalid input: playlistID are empty")
	}

	if len(trackIDs) == 0 {
		const query = `
		DELETE FROM playlist_tracks
		WHERE playlist_id = $1
	`
		_, err := s.db.ExecContext(ctx, query, playlistID)
		return err
	}

	const query = `
		DELETE FROM playlist_tracks
		WHERE playlist_id = $1 AND track_id = ANY($2)
	`
	_, err = s.db.ExecContext(ctx, query, playlistID, pq.Array(trackIDs))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
