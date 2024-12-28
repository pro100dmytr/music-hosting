package repository

import (
	"context"
	"database/sql"
	"music-hosting/internal/config"
	"music-hosting/internal/models"
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

func (s *UserStorage) Create(ctx context.Context, user *models.User) error {
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
		return err
	}

	user.ID = id
	return nil
}

func (s *UserStorage) Get(ctx context.Context, id int) (*models.User, error) {
	const query = `SELECT * FROM users WHERE id = $1`
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.Password,
	)
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
		user := &models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.Password,
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

func (s *UserStorage) Update(ctx context.Context, user *models.User, id int) error {
	const query = `UPDATE users SET login = $1, email = $2, password = $3 WHERE id = $4`
	result, err := s.db.ExecContext(ctx, query, user.Login, user.Email, user.Password, id)
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

//func (s *UserStorage) Save(ctx context.Context, playlistID int, track *models.Track) error {
//	query := `INSERT INTO tracks (name, artist, url, playlist_id) VALUES ($1, $2, $3, $4)`
//	_, err := s.db.ExecContext(ctx, query, track.Name, track.Artist, track.URL, playlistID)
//	return err
//}
