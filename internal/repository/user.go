package repository

import (
	"context"
	"database/sql"
	"errors"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) (*UserStorage, error) {
	return &UserStorage{db: db}, nil
}

func (s *UserStorage) Create(ctx context.Context, user *User) (int, error) {
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

func (s *UserStorage) Get(ctx context.Context, id int) (*User, error) {
	const query = `SELECT id, login, email FROM users WHERE id = $1`

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) Update(ctx context.Context, user *User, id int) error {
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

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM users WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetUsers(ctx context.Context, offset, limit int) ([]*User, error) {
	const query = `
        SELECT id, login, email FROM users OFFSET $1 LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (s *UserStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	const query = `SELECT id, login, password_hash, salt FROM users WHERE login = $1`
	user := &User{}
	err := s.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.Salt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) AddPlaylistsToUser(ctx context.Context, userID int, playlistID int) error {
	const query = `UPDATE playlists SET user_id = $1 WHERE id = $2`

	_, err := s.db.ExecContext(ctx, query, userID, playlistID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) RemovePlaylistsFromUser(ctx context.Context, userID int) error {
	const query = `DELETE FROM playlists WHERE user_id = $1`

	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) UpdatePlбaylistsForUser(ctx context.Context, userID int) error {
	const query = `UPDATE playlists SET updated_at = NOW() WHERE user_id = $1`

	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
