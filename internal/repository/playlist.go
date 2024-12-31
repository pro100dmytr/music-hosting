package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music-hosting/internal/config"
	"music-hosting/internal/models/repositorys"
	"music-hosting/pkg/database/postgresql"
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
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `INSERT INTO playlists (name, user_id, track_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err = s.db.QueryRowContext(
		ctx,
		query,
		playlist.Name,
		playlist.UserID,
		playlist.TrackID,
	).Scan(
		&playlist.ID,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PlaylistStorage) GetAll(ctx context.Context) ([]*repositorys.Playlist, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `SELECT id, name, user_id, track_id, created_at, updated_at FROM playlists`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []*repositorys.Playlist
	for rows.Next() {
		playlist := &repositorys.Playlist{}
		if err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.UserID,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

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

func (s *PlaylistStorage) Get(ctx context.Context, id int) (*repositorys.Playlist, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `SELECT id, name, user_id, track_id, created_at, updated_at FROM playlists WHERE id = $1`

	playlist := &repositorys.Playlist{}
	err = s.db.QueryRowContext(ctx, query, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.UserID,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return playlist, nil
}

func (s *PlaylistStorage) GetByUserID(ctx context.Context, userID int) ([]*repositorys.Playlist, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

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
			&playlist.ID,
			&playlist.Name,
			&playlist.UserID,
			&playlist.TrackID,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

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

func (s *PlaylistStorage) GetByName(ctx context.Context, name string) ([]*repositorys.Playlist, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `SELECT id, name, user_id, track_id, created_at, updated_at FROM playlists WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var playlists []*repositorys.Playlist
	for rows.Next() {
		playlist := &repositorys.Playlist{}
		if err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.UserID,
			&playlist.TrackID,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

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

func (s *PlaylistStorage) Update(ctx context.Context, id int, playlist *repositorys.Playlist) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `UPDATE playlists SET name = $1, user_id = $2, track_id = $3, updated_at = NOW() WHERE id = $4`

	result, err := s.db.ExecContext(ctx, query, playlist.Name, playlist.UserID, playlist.TrackID, id)
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

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
