package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/yourusername/shorty/internal/config"
)

func NewPgxPool(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgDB,
		cfg.PgSSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Error().Err(err).Msg("pgx connect")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_code VARCHAR(20) NOT NULL UNIQUE,
			original_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			clicks INT DEFAULT 0
		);
	`)
	if err != nil {
		pool.Close()
		log.Error().Err(err).Msg("create urls table")
		return nil, err
	}

	log.Info().Msg("pgx connected and urls table ensured")
	return pool, nil
}
