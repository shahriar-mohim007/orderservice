package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type PgxRepository struct {
	db *pgxpool.Pool
}

var (
	once       sync.Once
	repository *PgxRepository
)

func NewPgRepository(databaseUrl string) (*PgxRepository, error) {
	var onceErr error
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		config, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			onceErr = fmt.Errorf("invalid database URL: %w", err)
			log.Error().Err(err).Msg("Failed to parse database configuration")
			return
		}

		config.MaxConns = 100
		config.MinConns = 2
		config.MaxConnLifetime = 30 * time.Minute
		config.MaxConnIdleTime = 5 * time.Second
		config.HealthCheckPeriod = 1 * time.Minute

		db, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			onceErr = fmt.Errorf("failed to create connection pool: %w", err)
			log.Error().Err(err).Msg("Database Connection Error")
			return
		}

		if err = db.Ping(ctx); err != nil {
			onceErr = fmt.Errorf("failed to ping database: %w", err)
			log.Error().Err(err).Msg("Database Ping Error")
			db.Close()
			return
		}

		repository = &PgxRepository{db: db}
		log.Info().Msg("Database connection pool successfully initialized")
	})

	return repository, onceErr
}

func (repo *PgxRepository) Close() {
	if repo.db != nil {
		repo.db.Close()
		log.Info().Msg("Database connection pool closed")
	}
}
