package repository

import (
	"context"
	"database/sql"
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
	const query = `INSERT INTO tracks (name, artist, url, playlist_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := s.db.QueryRow(query, track.Name, track.Artist, track.URL, track.PlaylistID).Scan(&id)
	if err != nil {
		return err
	}

	track.ID = id
	return nil
}

func (s *TrackStorage) Get(ctx context.Context, id int) (*models.Track, error) {
	const query = `SELECT * FROM tracks WHERE id = $1`
	track := &models.Track{}
	err := s.db.QueryRow(query, id).Scan(
		&track.ID,
		&track.Name,
		&track.Artist,
		&track.URL,
		&track.PlaylistID,
	)
	if err != nil {
		return nil, err
	}

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
			&track.PlaylistID,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, rows.Err()
}

func (s *TrackStorage) Update(ctx context.Context, track *models.Track, id int) error {
	const query = `UPDATE tracks SET name = $1, artist = $2, url = $3, playlist_id = $4 WHERE id = $5`
	result, err := s.db.ExecContext(ctx, query, track.Name, track.Artist, track.URL, track.PlaylistID, id)
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
