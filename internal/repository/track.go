package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music-hosting/internal/config"
	"music-hosting/internal/database/postgresql"
	"music-hosting/internal/models/repositorys"
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
	const query = `INSERT INTO tracks (name, artist, url, likes, dislikes) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := s.db.QueryRowContext(ctx, query, track.Name, track.Artist, track.URL, track.Likes, track.Dislikes).Scan(&id)
	if err != nil {
		return 0, err
	}

	track.ID = id
	return id, nil
}

func (s *TrackStorage) Get(ctx context.Context, id int) (*repositorys.Track, error) {
	const query = `SELECT * FROM tracks WHERE id = $1`

	track := &repositorys.Track{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		track.ID,
		track.Name,
		track.Artist,
		track.URL,
		track.Likes,
		track.Dislikes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows // TODO: return  nil instead of sql.ErrNoRows
		}
		return nil, err
	}

	track.ID = id
	return track, nil
}

func (s *TrackStorage) GetAll(ctx context.Context) ([]*repositorys.Track, error) {
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
			track.ID,
			track.Name,
			track.Artist,
			track.URL,
			track.Likes,
			track.Dislikes,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}

// TODO: rename to GetByName
func (s *TrackStorage) GetForName(ctx context.Context, name string) ([]*repositorys.Track, error) {
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
			track.ID,
			track.Name,
			track.Artist,
			track.URL,
			track.Likes,
			track.Dislikes,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}

func (s *TrackStorage) GetForArtist(ctx context.Context, artist string) ([]*repositorys.Track, error) {
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
			track.ID,
			track.Name,
			track.Artist,
			track.URL,
			track.Likes,
			track.Dislikes,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}

// TODO: delete 'id' parameter. Take 'id' from 'track'
func (s *TrackStorage) Update(ctx context.Context, track *repositorys.Track, id int) error {
	// TODO: транзакция тут не нужна
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `UPDATE tracks SET name = $1, artist = $2, url = $3, likes = $4, dislikes = $5 WHERE id = $6`
	result, err := s.db.ExecContext(ctx, query, track.Name, track.Artist, track.URL, track.Artist, track.Dislikes, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// TODO: удали этот стейтмент. Он тут не нужен
	if n == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TrackStorage) Delete(ctx context.Context, id int) error {
	// TODO: транзакция тут не нужна
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
	const query = `SELECT id, name, artist, url, likes, dislikes FROM tracks OFFSET $1 LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*repositorys.Track
	for rows.Next() {
		track := &repositorys.Track{}
		if err := rows.Scan(
			track.ID,
			track.Name,
			track.Artist,
			track.URL,
			track.Likes,
			track.Dislikes,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}

// TODO: этот метод тебе не нужен. Это можно делать через Update
// В сервисе нужно получать влейлист, менять кол-во лайков и сохранять
func (s *TrackStorage) AddLike(ctx context.Context, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const checkQuery = `SELECT 1 FROM tracks WHERE id = $1`
	var exists bool
	err = s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to check if track exists: %w", err)
	}

	const query = `UPDATE tracks SET like = like + 1 WHERE id = $1`

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

// TODO: этот метод тебе не нужен. Это можно делать через Update.
// В сервисе нужно получать влейлист, менять кол-во лайков и сохранять
func (s *TrackStorage) RemoveLike(ctx context.Context, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const checkQuery = `SELECT 1 FROM tracks WHERE id = $1`
	var exists bool
	err = s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to check if track exists: %w", err)
	}

	const query = `UPDATE tracks SET like = like - 1 WHERE id = $1`

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

// TODO: этот метод тебе не нужен. Это можно делать через Update.
// В сервисе нужно получать влейлист, менять кол-во лайков и сохранять
func (s *TrackStorage) AddDislike(ctx context.Context, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const checkQuery = `SELECT 1 FROM tracks WHERE id = $1`
	var exists bool
	err = s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to check if track exists: %w", err)
	}

	const query = `UPDATE tracks SET dislike = dislike + 1 WHERE id = $1`

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

// TODO: этот метод тебе не нужен. Это можно делать через Update.
// В сервисе нужно получать влейлист, менять кол-во лайков и сохранять
func (s *TrackStorage) RemoveDislike(ctx context.Context, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const checkQuery = `SELECT 1 FROM tracks WHERE id = $1`
	var exists bool
	err = s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return fmt.Errorf("failed to check if track exists: %w", err)
	}

	const query = `UPDATE tracks SET dislike = dislike - 1 WHERE id = $1`

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
