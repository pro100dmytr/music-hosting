package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music-hosting/internal/config"
	"music-hosting/internal/database/postgresql"
	"music-hosting/internal/models/repositorys"

	"github.com/lib/pq"
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

func (s *UserStorage) Create(ctx context.Context, user *repositorys.User) (int, error) {
	const query = `INSERT INTO users (login, email, password_hash, salt) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Login,
		user.Email,
		user.Password,
		user.Salt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserStorage) Get(ctx context.Context, id int) (*repositorys.User, error) {
	const query = `SELECT id, login, email FROM users WHERE id = $1`

	user := &repositorys.User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		user.ID,
		user.Login,
		user.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows //  TODO: return nil instead of sql.ErrNoRows
		}
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]*repositorys.User, error) {
	const query = `SELECT id, login, email FROM users`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*repositorys.User
	for rows.Next() {
		user := &repositorys.User{}
		if err := rows.Scan(
			user.ID,
			user.Login,
			user.Email,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (s *UserStorage) Update(ctx context.Context, user *repositorys.User, id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()
	const query = `UPDATE users SET login = $1, email = $2, password_hash = $3, salt = $4 WHERE id = $5`

	result, err := s.db.ExecContext(
		ctx,
		query,
		user.Login,
		user.Email,
		user.Password,
		user.Salt,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// TODO:  delete statement
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *UserStorage) Delete(ctx context.Context, id int) error {
	// TODO: транзакция тут не нужна
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()
	const query = `DELETE FROM users WHERE id = $1`

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

func (s *UserStorage) GetUsers(ctx context.Context, offset, limit int) ([]*repositorys.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	const query = `
        SELECT id, login, email FROM users OFFSET $1 LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*repositorys.User
	for rows.Next() {
		user := &repositorys.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return users, rows.Err()
}

func (s *UserStorage) GetUserByLogin(ctx context.Context, login string) (*repositorys.User, error) {
	user := &repositorys.User{}

	const query = `SELECT id, login, password_hash, salt FROM users WHERE login = $1`
	err := s.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.Salt)
	if err != nil {
		// TODO: стейтмент тут не нужен
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

// TODO: сделай чтобы работало
func (s *UserStorage) AddPlaylistsToUser(ctx context.Context, userID int, playlistIDs []int) error {
	const query = `
		INSERT INTO playlist_tracks (playlist_id, track_id)
		SELECT unnest($2::int[]), track_id
		FROM playlists
		WHERE user_id = $1 AND playlist_id = ANY($2)
		ON CONFLICT DO NOTHING
	`

	_, err := s.db.ExecContext(ctx, query, userID, pq.Array(playlistIDs))
	if err != nil {
		return err
	}

	return nil
}

// TODO: сделай чтобы работало
func (s *UserStorage) RemovePlaylistsFromUser(ctx context.Context, userID int, playlistIDs []int) error {
	const query = `
		DELETE FROM playlist_tracks
		WHERE playlist_id = ANY($2)
		AND EXISTS (
			SELECT 1 FROM playlists
			WHERE user_id = $1 AND playlist_id = playlist_tracks.playlist_id
		)
	`

	_, err := s.db.ExecContext(ctx, query, userID, pq.Array(playlistIDs))
	if err != nil {
		return err
	}

	return nil
}

// TODO: этот метод должен быть в PlaylistRepository и должен возвращаться []Playlist
func (s *UserStorage) GetPlaylistsForUser(ctx context.Context, userID int) ([]int, error) {
	const query = `SELECT id FROM playlists WHERE user_id = $1`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlistIDs []int
	for rows.Next() {
		var playlistID int
		if err := rows.Scan(&playlistID); err != nil {
			return nil, err
		}
		playlistIDs = append(playlistIDs, playlistID)
	}

	return playlistIDs, nil
}

// TODO: сделай чтобы работало
func (s *UserStorage) UpdatePlaylistsForUser(ctx context.Context, userID int, playlistIDs []int) error {
	const query = `
		UPDATE playlists
		SET updated_at = NOW()
		WHERE user_id = $1 AND playlist_id = ANY($2)
	`

	_, err := s.db.ExecContext(ctx, query, userID, pq.Array(playlistIDs))
	if err != nil {
		return err
	}

	return nil
}
