package repository

import (
	"context"
	"database/sql"
	"fmt"
	"music-hosting/internal/config"
	"music-hosting/internal/models"
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

func (s *TrackStorage) Create(ctx context.Context, track *models.Track) error {
	const query = `INSERT INTO tracks (name, artist, url) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.db.QueryRowContext(ctx, query, track.Name, track.Artist, track.URL).Scan(&id)
	if err != nil {
		return err
	}

	track.ID = id
	return nil
}

func (s *TrackStorage) Get(ctx context.Context, id int) (*models.Track, error) {
	const query = `SELECT * FROM tracks WHERE id = $1`
	track := &models.Track{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&track.ID,
		&track.Name,
		&track.Artist,
		&track.URL,
	)
	if err != nil {
		return nil, err
	}

	track.ID = id
	return track, nil
}

func (s *TrackStorage) GetAll(ctx context.Context) ([]*models.Track, error) {
	const query = `SELECT * FROM tracks`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}
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
	return tracks, rows.Err()
}

func (s *TrackStorage) GetForName(ctx context.Context, name string) ([]*models.Track, error) {
	const query = `SELECT * FROM tracks WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}
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
	return tracks, rows.Err()
}

func (s *TrackStorage) GetForArtist(ctx context.Context, artist string) ([]*models.Track, error) {
	const query = `SELECT * FROM tracks WHERE artist = $1`
	rows, err := s.db.QueryContext(ctx, query, artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}
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
	return tracks, rows.Err()
}

func (s *TrackStorage) Update(ctx context.Context, track *models.Track, id int) error {
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

	return nil
}

func (s *TrackStorage) Delete(ctx context.Context, id int) error {
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

	return nil
}

func (s *TrackStorage) GetTracks(ctx context.Context, limit, offset int) ([]*models.Track, error) {
	const query = `SELECT id, name, artist, url FROM tracks LIMIT $1 OFFSET $2`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracks: %w", err)
	}
	defer rows.Close()

	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}

		if err := rows.Scan(&track.ID, &track.Name, &track.Artist, &track.URL); err != nil {
			return nil, fmt.Errorf("failed to scan track row: %w", err)
		}

		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tracks, nil
}
