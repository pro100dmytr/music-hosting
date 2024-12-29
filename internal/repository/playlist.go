package repository

import (
	"context"
	"database/sql"
	"music-hosting/internal/config"
	"music-hosting/internal/models"
	"music-hosting/pkg/database/postgresql"
	"music-hosting/pkg/utils/convertutils"
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
	const query = `INSERT INTO playlists (name, user_id, track_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	IdArray := convertutils.IntSliceConvertIntoString(playlist.TrackID)

	err := s.db.QueryRowContext(
		ctx,
		query,
		playlist.Name,
		playlist.UserID,
		IdArray,
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

func (s *PlaylistStorage) GetAll(ctx context.Context) ([]*models.Playlist, error) {
	const query = `SELECT id, name, user_id, track_id, created_at, updated_at FROM playlists`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []*models.Playlist
	for rows.Next() {
		var trackIDRaw []byte

		playlist := &models.Playlist{}
		if err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.UserID,
			&trackIDRaw,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

		playlist.TrackID, err = convertutils.StringConvertIntoIntSlice(string(trackIDRaw))
		if err != nil {
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
	const query = `SELECT id, name, user_id, track_id, created_at, updated_at FROM playlists WHERE id = $1`
	var trackIDRaw []byte

	playlist := &models.Playlist{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.UserID,
		&trackIDRaw,
		&playlist.CreatedAt,
		&playlist.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	playlist.TrackID, err = convertutils.StringConvertIntoIntSlice(string(trackIDRaw))
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistStorage) GetByUserID(ctx context.Context, userID int) ([]*models.Playlist, error) {
	const query = `SELECT * FROM playlists WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var playlists []*models.Playlist
	for rows.Next() {
		var trackIDRaw []byte
		playlist := &models.Playlist{}
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

		playlist.TrackID, err = convertutils.StringConvertIntoIntSlice(string(trackIDRaw))
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func (s *PlaylistStorage) GetByName(ctx context.Context, name string) ([]*models.Playlist, error) {
	const query = `SELECT id, name, user_id, track_id, created_at, updated_at FROM playlists WHERE name = $1`
	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var playlists []*models.Playlist
	for rows.Next() {
		var trackIDRaw string
		playlist := &models.Playlist{}
		if err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.UserID,
			&trackIDRaw,
			&playlist.CreatedAt,
			&playlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

		playlist.TrackID, err = convertutils.StringConvertIntoIntSlice(trackIDRaw)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}

func (s *PlaylistStorage) Update(ctx context.Context, id int, playlist *models.Playlist) error {
	const checkTrackQuery = `SELECT COUNT(*) FROM tracks WHERE id = $1`
	for _, trackID := range playlist.TrackID {
		var count int
		err := s.db.QueryRowContext(ctx, checkTrackQuery, trackID).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			return sql.ErrNoRows
		}
	}

	const query = `UPDATE playlists SET name = $1, user_id = $2, track_id = $3, updated_at = NOW() WHERE id = $4`
	trackIDString := convertutils.IntSliceConvertIntoString(playlist.TrackID)

	_, err := s.db.ExecContext(ctx, query, playlist.Name, playlist.UserID, trackIDString, id)
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
