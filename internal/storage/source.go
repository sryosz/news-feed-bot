package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"news-feed-bot/internal/model"
	"os"
)

type SourcePostgresStorage struct {
	db *sql.DB
}

func NewSourceStorage(log *slog.Logger) (*SourcePostgresStorage, error) {
	const op = "storage.postgres.New"

	log.Info("connecting to db | Source storage")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("connected to db successfully")

	return &SourcePostgresStorage{db: db}, nil
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	const op = "storage.source.Sources"

	stmt, err := s.db.QueryContext(ctx, "SELECT * FROM sources")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var sources []model.Source

	for stmt.Next() {
		var source model.Source
		if err := stmt.Scan(&source.ID, &source.Name, &source.FeedURL, &source.CreatedAt, &source.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (s *SourcePostgresStorage) SourceById(ctx context.Context, id int64) (*model.Source, error) {
	const op = "storage.source.SourceById"

	stmt, err := s.db.Prepare("SELECT * FROM sources WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var source model.Source

	err = stmt.QueryRowContext(ctx, id).
		Scan(
			&source.ID,
			&source.Name,
			&source.FeedURL,
			&source.CreatedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &source, nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	const op = "storage.source.Add"

	stmt, err := s.db.Prepare("INSERT INTO sources(name, feed_url, created_at) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64

	err = stmt.QueryRowContext(ctx, source.Name, source.FeedURL, source.CreatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	const op = "storage.source.Delete"

	_, err := s.db.ExecContext(ctx, "DELETE FROM sources WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
