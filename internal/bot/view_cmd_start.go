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
		// TODO: remake from hardcoded list of available commands to for loop with joining
		msgTxt := fmt.Sprintf("Greetings!" +
			"\n\nThis bot is designed for scheduled publication\nof articles from specified sources" +
			"\n\nType /commands	 to see available commands")

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, msgTxt)); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}
