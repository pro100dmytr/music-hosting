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

type TrackStorage struct {
	db *sql.DB
}

func (s *TrackStorage) Close() error {
	return postgresql.CloseConn(s.db)
}

func NewTrackStorage(cfg *config.Config) (*TrackStorage, error) {
	db, err := postgresql.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &TrackStorage{db: db}, nil
}

func (s *TrackStorage) Create(ctx context.Context, track *repositorys.Track) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `INSERT INTO tracks (name, artist, url) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err = s.db.QueryRowContext(ctx, query, track.Name, track.Artist, track.URL).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	track.ID = id
	return id, nil
}

func (s *TrackStorage) Get(ctx context.Context, id int) (*repositorys.Track, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `SELECT * FROM tracks WHERE id = $1`

	track := &repositorys.Track{}
	err = s.db.QueryRowContext(ctx, query, id).Scan(
		&track.ID,
		&track.Name,
		&track.Artist,
		&track.URL,
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

	track.ID = id
	return track, nil
}

func (s *TrackStorage) GetAll(ctx context.Context) ([]*repositorys.Track, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()
	const query = `SELECT * FROM tracks`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*repositorys.Track
	for rows.Next() {
		track := &repositorys.Track{}
		if err := rows.Scan(
			&track.ID,
			&track.Name,
			&track.Artist,
			&track.URL,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tracks, rows.Err()
}

func (s *TrackStorage) GetForName(ctx context.Context, name string) ([]*repositorys.Track, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `SELECT * FROM tracks WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*repositorys.Track
	for rows.Next() {
		track := &repositorys.Track{}
		if err := rows.Scan(
			&track.ID,
			&track.Name,
			&track.Artist,
			&track.URL,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tracks, rows.Err()
}

func (s *TrackStorage) GetForArtist(ctx context.Context, artist string) ([]*repositorys.Track, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `SELECT * FROM tracks WHERE artist = $1`
	rows, err := s.db.QueryContext(ctx, query, artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*repositorys.Track
	for rows.Next() {
		track := &repositorys.Track{}
		if err := rows.Scan(
			&track.ID,
			&track.Name,
			&track.Artist,
			&track.URL,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tracks, rows.Err()
}

func (s *TrackStorage) Update(ctx context.Context, track *repositorys.Track, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `UPDATE tracks SET name = $1, artist = $2, url = $3 WHERE id = $4`
	result, err := s.db.ExecContext(ctx, query, track.Name, track.Artist, track.URL, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TrackStorage) Delete(ctx context.Context, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `DELETE FROM tracks WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TrackStorage) GetTracks(ctx context.Context, offset, limit int) ([]*repositorys.Track, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `SELECT id, name, artist, url FROM tracks OFFSET $1 LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*repositorys.Track
	for rows.Next() {
		track := &repositorys.Track{}
		if err := rows.Scan(
			&track.ID,
			&track.Name,
			&track.Artist,
			&track.URL,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tracks, rows.Err()
}
