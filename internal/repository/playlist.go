package repository

import (
	"context"
	"database/sql"
	"music-hosting/internal/config"
	"music-hosting/internal/models"
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

func (s *PlaylistStorage) Create(ctx context.Context, playlist *models.Playlist) error {
	const query = `INSERT INTO playlists (name, user_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := s.db.QueryRow(query, playlist.Name, playlist.UserID).Scan(&playlist.ID, &playlist.CreatedAt, &playlist.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistStorage) GetAll(ctx context.Context) ([]*models.Playlist, error) {
	const query = `SELECT id, name, user_id, created_at, updated_at FROM playlists`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []*models.Playlist
	for rows.Next() {
		playlist := &models.Playlist{}
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

	return playlists, nil
}

func (s *PlaylistStorage) Get(ctx context.Context, id int) (*models.Playlist, error) {
	const query = `SELECT id, name, user_id, created_at, updated_at FROM playlists WHERE id = $1`
	playlist := &models.Playlist{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.UserID,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistStorage) Update(ctx context.Context, id int, playlist *models.Playlist) error {
	const query = `UPDATE playlists SET name = $1, user_id = $2, updated_at = NOW() WHERE id = $3`
	_, err := s.db.ExecContext(ctx, query, playlist.Name, playlist.UserID, id)
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
