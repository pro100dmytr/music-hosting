package postgresql

import (
	"database/sql"
	"fmt"
	"music-hosting/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func OpenConnection(cfg *config.DBConfig) (*sql.DB, error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Password == "" || cfg.DBName == "" || cfg.SSLMode == "" {
		return nil, fmt.Errorf("incomplete storage configuration: host=%s, port=%s, user=%s, dbname=%s, sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.DBName,
			cfg.SSLMode,
		)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping storage: %w", err)
	}

	err = goose.Up(db, "db/migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}
