package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
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

	return track, nil
}

func (s *TrackStorage) GetTracks(ctx context.Context, name, artist string, playlistID, offset, limit int) ([]*Track, error) {
	baseQuery := `SELECT t.id, t.name, t.artist, t.url, t.likes, t.dislikes FROM tracks t`
	var conditions []string
	var args []interface{}

	if name != "" {
		conditions = append(conditions, "t.name = $"+strconv.Itoa(len(args)+1))
		args = append(args, name)
	}

	if artist != "" {
		conditions = append(conditions, "t.artist = $"+strconv.Itoa(len(args)+1))
		args = append(args, artist)
	}

	if playlistID > 0 {
		baseQuery += ` JOIN playlist_tracks pt ON pt.track_id = t.id`
		conditions = append(conditions, "pt.playlist_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, playlistID)
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	if limit > 0 {
		baseQuery += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, limit)
	}
	if offset >= 0 {
		baseQuery += " OFFSET $" + strconv.Itoa(len(args)+1)
		args = append(args, offset)
	}

	rows, err := s.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*Track
	for rows.Next() {
		track := &Track{}
		if err := rows.Scan(&track.ID, &track.Name, &track.Artist, &track.URL, &track.Likes, &track.Dislikes); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, rows.Err()
}

func (s *TrackStorage) Update(ctx context.Context, track *Track) error {
	const query = `UPDATE tracks SET name = $1, artist = $2, url = $3, likes = $4, dislikes = $5 WHERE id = $6`
	_, err := s.db.ExecContext(
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

	return nil
}

func (s *TrackStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM tracks WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
