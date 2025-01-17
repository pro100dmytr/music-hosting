package repository

import (
	"context"
	"database/sql"
	"errors"
)

type TrackStorage struct {
	db *sql.DB
}

func NewTrackStorage(db *sql.DB) (*TrackStorage, error) {
	return &TrackStorage{db: db}, nil
}

func (s *TrackStorage) Create(ctx context.Context, track *Track) (int, error) {
	const query = `INSERT INTO tracks (name, artist, url, likes, dislikes) VALUES ($1, $2, $3, $4, $5)`
	var id int
	_, err := s.db.QueryContext(ctx, query, track.Name, track.Artist, track.URL, track.Likes, track.Dislikes)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TrackStorage) Get(ctx context.Context, id int) (*Track, error) {
	const query = `SELECT * FROM tracks WHERE id = $1`

	track := &Track{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&track.ID,
		&track.Name,
		&track.Artist,
		&track.URL,
		&track.Likes,
		&track.Dislikes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	track.ID = id
	return track, nil
}

func (s *TrackStorage) GetAll(ctx context.Context) ([]*Track, error) {
	const query = `SELECT * FROM tracks`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*Track
	for rows.Next() {
		track := &Track{}
		if err := rows.Scan(
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

	return tracks, rows.Err()
}

func (s *TrackStorage) GetByName(ctx context.Context, name string) ([]*Track, error) {
	const query = `SELECT * FROM tracks WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*Track
	for rows.Next() {
		track := &Track{}
		if err := rows.Scan(
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

	return tracks, rows.Err()
}

func (s *TrackStorage) GetByArtist(ctx context.Context, artist string) ([]*Track, error) {
	const query = `SELECT * FROM tracks WHERE artist = $1`
	rows, err := s.db.QueryContext(ctx, query, artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*Track
	for rows.Next() {
		track := &Track{}
		if err := rows.Scan(
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

	return tracks, rows.Err()
}

func (s *TrackStorage) Update(ctx context.Context, track *Track) error {
	const query = `UPDATE tracks SET name = $1, artist = $2, url = $3, likes = $4, dislikes = $5 WHERE id = $6`
	result, err := s.db.ExecContext(
		ctx,
		query,
		track.Name,
		track.Artist,
		track.URL,
		track.Likes,
		track.Dislikes,
		track.ID,
	)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
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

func (s *TrackStorage) GetTracks(ctx context.Context, offset, limit int) ([]*Track, error) {
	const query = `SELECT id, name, artist, url, likes, dislikes FROM tracks OFFSET $1 LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*Track
	for rows.Next() {
		track := &Track{}
		if err := rows.Scan(
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

	return tracks, rows.Err()
}

func (s *TrackStorage) GetTracksByPlaylistID(ctx context.Context, playlistID int) ([]*Track, error) {
	const query = `SELECT track_id FROM playlist_tracks WHERE playlist_id = $1`

	rows, err := s.db.QueryContext(ctx, query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracksID []*Track
	for rows.Next() {
		var trackID *Track
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
