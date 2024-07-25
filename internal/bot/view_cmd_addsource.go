package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/model"
	"time"
)

type SourceStorage interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

func ViewCmdAddSource(storage SourceStorage) botkit.ViewFunc {
	const op = "bot.ViewCmdAddSource"

	type addSourceArgs struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		source := model.Source{
			Name:      args.Name,
			FeedURL:   args.URL,
			CreatedAt: time.Now().UTC(),
		}

		sourceID, err := storage.Add(ctx, source)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		var (
			msgText = fmt.Sprintf(
				"Source was added with id: `%d`\\."+
					" Use this id for managing this source\\.", sourceID)
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}
