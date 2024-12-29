package repository

import (
	"context"
	"database/sql"
	"fmt"
	"music-hosting/internal/config"
	"music-hosting/internal/models"
	"music-hosting/pkg/database/postgresql"
	"music-hosting/pkg/utils/convertutils"
)

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Close() error {
	return postgresql.CloseConn(s.db)
}

func NewUserStorage(cfg *config.Config) (*UserStorage, error) {
	db, err := postgresql.OpenConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &UserStorage{db: db}, nil
}

func (s *UserStorage) Create(ctx context.Context, user *models.User) (int, error) {
	const query = `INSERT INTO users (login, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Login,
		user.Email,
		user.Password,
	).Scan(&id)

	if err != nil {
		return user.ID, err
	}

	user.ID = id
	return user.ID, nil
}

func (s *UserStorage) Get(ctx context.Context, id int) (*models.User, error) {
	const query = `SELECT * FROM users WHERE id = $1`
	var userIDRaw []byte
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.Password,
		&userIDRaw,
	)
	if err != nil {
		return nil, err
	}

	user.PlaylistID, err = convertutils.StringConvertIntoIntSlice(string(userIDRaw))
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]*models.User, error) {
	const query = `SELECT * FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var playlistIDRaw []byte

		user := &models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.Password,
			&playlistIDRaw,
		); err != nil {
			return nil, err
		}

		user.PlaylistID, err = convertutils.StringConvertIntoIntSlice(string(playlistIDRaw))
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) Update(ctx context.Context, user *models.User, id int) error {
	const checkPlaylistQuery = `SELECT COUNT(*) FROM playlists WHERE id = $1`
	for _, playlistID := range user.PlaylistID {
		var count int
		err := s.db.QueryRowContext(ctx, checkPlaylistQuery, playlistID).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			return sql.ErrNoRows
		}
	}

	const query = `UPDATE users SET login = $1, email = $2, password = $3, playlist_id = $4 WHERE id = $5`
	playlistIDString := convertutils.IntSliceConvertIntoString(user.PlaylistID)

	result, err := s.db.ExecContext(ctx, query, user.Login, user.Email, user.Password, playlistIDString, id)
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

	user.ID = id
	return nil
}

func (s *UserStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM users WHERE id = $1`
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

func (s *UserStorage) GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	const query = `SELECT id, login, email, playlist_id FROM users LIMIT $1 OFFSET $2`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var playlistIDRaw []byte

		if err := rows.Scan(&user.ID, &user.Login, &user.Email, &playlistIDRaw); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user.PlaylistID, err = convertutils.StringConvertIntoIntSlice(string(playlistIDRaw))
		if err != nil {
			return nil, fmt.Errorf("failed to convert playlist_id: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}
