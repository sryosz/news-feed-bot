package middleware

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
)

func AdminOnly(channelID int64, next botkit.ViewFunc) botkit.ViewFunc {
	const op = "bot.middleware.AdminOnly"

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		admins, err := bot.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelID,
				},
			},
		)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
				return next(ctx, bot, update)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"You have no permissions"),
		); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}
