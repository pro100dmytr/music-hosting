package repository

import (
	"context"
	"database/sql"
	"errors"
	"music-hosting/internal/config"
	"music-hosting/internal/models/repositorys"
	"music-hosting/pkg/database/postgresql"
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
	const query = `INSERT INTO users (login, email, password, playlist_id) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Login,
		user.Email,
		user.Password,
		user.PlaylistID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserStorage) Get(ctx context.Context, id int) (*repositorys.User, error) {
	const query = `
        SELECT id, login, email, password, playlist_id 
        FROM users 
        WHERE id = $1`

	user := &repositorys.User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.Password,
		&user.PlaylistID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]*repositorys.User, error) {
	const query = `
        SELECT id, login, email, password, playlist_id 
        FROM users`

	rows, err := s.db.QueryContext(ctx, query)
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
			&user.Password,
			&user.PlaylistID,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (s *UserStorage) Update(ctx context.Context, user *repositorys.User, id int) error {
	const query = `
        UPDATE users 
        SET login = $1, email = $2, password = $3, playlist_id = $4 
        WHERE id = $5`

	result, err := s.db.ExecContext(
		ctx,
		query,
		user.Login,
		user.Email,
		user.Password,
		user.PlaylistID,
		id,
	)
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

	return nil
}

func (s *UserStorage) Delete(ctx context.Context, id int) error {
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

	return nil
}

func (s *UserStorage) GetUsers(ctx context.Context, offset, limit int) ([]*repositorys.User, error) {
	const query = `
        SELECT id, login, email, password, playlist_id 
        FROM users 
        OFFSET $1 
        LIMIT $2`

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
			&user.Password,
			&user.PlaylistID,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) GetUserByLogin(ctx context.Context, login string) (*repositorys.User, error) {
	user := &repositorys.User{}

	query := `SELECT id, login, password FROM users WHERE login = $1`
	err := s.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}
