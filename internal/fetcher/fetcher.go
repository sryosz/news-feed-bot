package fetcher

import (
	"context"
	"fmt"
	"log/slog"
	"news-feed-bot/internal/model"
	"news-feed-bot/internal/source"
	"strings"
	"sync"
	"time"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articles       ArticleStorage
	sources        SourceProvider
	fetchInterval  time.Duration
	filterKeywords []string
	log            *slog.Logger
}

func New(
	articleStorage ArticleStorage,
	sourceProvider SourceProvider,
	fetchInterval time.Duration,
	filterKeywords []string,
	log *slog.Logger,
) *Fetcher {
	return &Fetcher{
		articles:       articleStorage,
		sources:        sourceProvider,
		log:            log,
		fetchInterval:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}

func (f *Fetcher) Start(ctx context.Context) error {
	const op = "fetcher.Start"

	f.log.Info("fetcher was started successfully")

	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: %w", op, ctx.Err())
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	const op = "fetcher.Fetch"

	f.log.Info("fetching sources")

	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)

		rssSource := source.New(src)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				//TODO: log error
				return
			}

			if err := f.processItems(ctx, source, items); err != nil {
				//TODO: log error
				return
			}
		}(rssSource)
	}

	wg.Wait()

	return nil
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.filterItem(item) {
			continue
		}

		if err := f.articles.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) filterItem(item model.Item) bool {
	for _, keyword := range f.filterKeywords {
		titleContainsKeyword := strings.Contains(strings.ToLower(item.Title), keyword)

		if titleContainsKeyword {
			return true
		}
	}

	return false
}
