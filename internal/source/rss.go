package source

import (
	"context"
	"fmt"
	"github.com/SlyMarbo/rss"
	"news-feed-bot/internal/model"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName string
}

func New(m model.Source) RSSSource {
	return RSSSource{
		URL:        m.FeedURL,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}

func (r RSSSource) Fetch(ctx context.Context) ([]model.Item, error) {
	const op = "source.rss.Fetch"
	// TODO: setup logger and log everything
	feed, err := r.loadFeed(ctx, r.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var items []model.Item

	for _, item := range feed.Items {
		items = append(items, model.Item{
			Title:      item.Title,
			Categories: item.Categories,
			Link:       item.Link,
			Date:       item.Date,
			Summary:    item.Summary,
			SourceName: r.SourceName,
		})
	}

	return items, nil
}

func (r RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	const op = "source.rss.loadFeed"
	// TODO: setup logger and log everything
	var (
		feedCh  = make(chan *rss.Feed)
		errorCh = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errorCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	case err := <-errorCh:
		return nil, fmt.Errorf("%s: %w", op, err)
	case feed := <-feedCh:
		return feed, nil
	}
}

func (r RSSSource) ID() int64 {
	return r.SourceID
}

func (r RSSSource) Name() string {
	return r.SourceName
}
