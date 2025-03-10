package notifier

import (
	"context"
	"fmt"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log/slog"
	"net/http"
	"news-feed-bot/internal/botkit/markup"
	"news-feed-bot/internal/model"
	"regexp"
	"strings"
	"time"
)

type ArticleProvider interface {
	NotPostedArticles(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkAsPosted(ctx context.Context, id int64) error
}

type Summarizer interface {
	Summarize(ctx context.Context, text string) (string, error)
}

type Notifier struct {
	articles         ArticleProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	articleRelevance time.Duration
	channelID        int64
	log              *slog.Logger
}

func New(
	articleProvider ArticleProvider,
	summarizer Summarizer,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	articleRelevance time.Duration,
	channelID int64,
	log *slog.Logger,
) *Notifier {
	return &Notifier{
		articles:         articleProvider,
		summarizer:       summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		articleRelevance: articleRelevance,
		channelID:        channelID,
		log:              log,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	const op = "notifier.Start"

	n.log.Info("notifier was started successfully")

	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case <-ctx.Done():
			return fmt.Errorf("%s: %w", op, ctx.Err())
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	const op = "notifier.SelectAndSendArticle"

	n.log.Info("selecting and sending article")

	topArticles, err := n.articles.NotPostedArticles(ctx, time.Now().Add(-n.articleRelevance), 1)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(topArticles) == 0 {
		return nil
	}

	article := topArticles[0]

	summary, err := n.extractSummary(ctx, article)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := n.sendArticle(article, summary); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return n.articles.MarkAsPosted(ctx, article.ID)
}

func (n *Notifier) extractSummary(ctx context.Context, article model.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}

	summary, err := n.summarizer.Summarize(ctx, cleanExtraEmptyLines(doc.TextContent))
	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil
}

func cleanExtraEmptyLines(text string) string {
	return regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n")
}

func (n *Notifier) sendArticle(article model.Article, summary string) error {
	const msgFormat = "*%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(n.channelID, fmt.Sprintf(
		msgFormat,
		markup.EscapeForMarkdown(article.Title),
		markup.EscapeForMarkdown(summary),
		markup.EscapeForMarkdown(article.Link),
	))
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
