package model

import (
	"database/sql"
	"time"
)

type Item struct {
	Title      string
	Categories []string
	Link       string
	Date       time.Time
	Summary    string
	SourceName string
}

type Source struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Article struct {
	ID          int64        `db:"id"`
	SourceID    int64        `db:"source_id"`
	Title       string       `db:"title"`
	Link        string       `db:"link"`
	Summary     string       `db:"summary"`
	PublishedAt time.Time    `db:"published_at"`
	CreatedAt   time.Time    `db:"created_at"`
	PostedAt    sql.NullTime `db:"posted_at"`
}
