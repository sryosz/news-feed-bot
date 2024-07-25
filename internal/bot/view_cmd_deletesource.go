package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
)

type SourceRemover interface {
	Delete(ctx context.Context, id int64) error
}

func ViewCmdDeleteSource(storage SourceRemover) botkit.ViewFunc {
	const op = "bot.ViewCmdAddSource"

	type deleteSourceArgs struct {
		ID int64 `json:"id"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[deleteSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = storage.Delete(ctx, args.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		var (
			msgText = fmt.Sprintf("Source was deleted with id: `%d`", args.ID)
			reply   = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}
