package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/botkit/markup"
	"news-feed-bot/internal/model"
	"strings"
)

type SourceLister interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

func ViewCmdListSources(lister SourceLister) botkit.ViewFunc {
	const op = "bot.ViewCmdListSources"

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := lister.Sources(ctx)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		var sourceInfos []string

		for _, source := range sources {
			sourceInfos = append(sourceInfos, formatSource(source))
		}

		msgText := fmt.Sprintf(
			"List of sources \\(total %d\\):\n\n%s",
			len(sources),
			strings.Join(sourceInfos, "\n\n"),
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}

func formatSource(source model.Source) string {
	return fmt.Sprintf(
		"⚪ *%s*\nID: `%d`\nFeed URL: %s",
		markup.EscapeForMarkdown(source.Name),
		source.ID,
		markup.EscapeForMarkdown(source.FeedURL),
	)
}
