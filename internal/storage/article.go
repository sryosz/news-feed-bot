package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"news-feed-bot/internal/model"
	"os"
	"time"
)

type ArticlePostgresStorage struct {
	db *sql.DB
}

func NewArticleStorage(log *slog.Logger) (*ArticlePostgresStorage, error) {
	const op = "storage.article.New"

	log.Info("connecting to db | Article storage")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("connected to db successfully")

	return &ArticlePostgresStorage{db: db}, nil
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	const op = "storage.article.Store"

	stmt, err := s.db.Prepare(`INSERT INTO articles
	                (source_id, title, link, summary, published_at)
					VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.ExecContext(
		ctx,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *ArticlePostgresStorage) NotPostedArticles(
	ctx context.Context,
	since time.Time,
	limit uint64,
) ([]model.Article, error) {
	const op = "storage.article.NotPostedArticles"

	stmt, err := s.db.Prepare(`SELECT * FROM articles
         WHERE posted_at IS NULL AND published_at > $1::timestamp
         ORDER BY  published_at DESC LIMIT $2`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var articles []model.Article

	rows, err := stmt.QueryContext(ctx, since, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := rows.Scan(&articles); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return articles, nil
}

func (s *ArticlePostgresStorage) MarkAsPosted(ctx context.Context, id int64) error {
	const op = "storage.article.MarkAsPosted"

	stmt, err := s.db.Prepare("UPDATE articles SET posted_at = $1::timestamp WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.ExecContext(ctx,
		time.Now().UTC().Format(time.RFC3339),
		id,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
