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
		msgTxt := fmt.Sprintf("Available commands:\n1. /addsource - adding new source\n2. /listsources - listing all sources\n")

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, msgTxt)); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}
