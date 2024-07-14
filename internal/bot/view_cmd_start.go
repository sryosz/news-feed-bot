package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
)

func ViewCmdStart() botkit.ViewFunc {
	const op = "bot.view_cmd_start.ViewCmdStart"

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Hello")); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}
