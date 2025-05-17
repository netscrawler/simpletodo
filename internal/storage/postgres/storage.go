package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	log     *zap.Logger
	Pool    *pgxpool.Pool
	Builder *squirrel.StatementBuilderType
}

func New(log *zap.Logger) *Storage {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Storage{log: log, Builder: &builder}
}

func (s *Storage) Init(ctx context.Context, connectionString string) error {
	const op = "storage.postgres.Init"

	var err error

	s.Pool, err = pgxpool.Connect(ctx, connectionString)
	if err != nil {
		s.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	if err := s.Pool.Ping(ctx); err != nil {
		s.log.Error(op, zap.Error(err))

		return fmt.Errorf("%s : %w", op, err)
	}

	s.log.Info(op + " : successfully connected")

	return nil
}

func (s *Storage) Close() {
	s.Pool.Close()
}
