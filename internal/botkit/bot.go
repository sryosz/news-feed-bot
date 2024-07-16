package botkit

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"time"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	CmdViews map[string]ViewFunc
	log      *slog.Logger
}

type ViewFunc func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(log *slog.Logger, api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		api: api,
		log: log,
	}
}

func (b *Bot) Run(ctx context.Context) error {
	const op = "botkit.Run"

	b.log.Info("bot was started successfully")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return fmt.Errorf("%s: %w", op, ctx.Err())
		}
	}
}

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.CmdViews == nil {
		b.CmdViews = make(map[string]ViewFunc)
	}

	b.CmdViews[cmd] = view
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			b.log.Info("%w", p)
		}
	}()

	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	var view ViewFunc

	cmd := update.Message.Command()

	cmdView, ok := b.CmdViews[cmd]
	if !ok {
		return
	}

	view = cmdView

	if err := view(ctx, b.api, update); err != nil {
		b.log.Info("%w", err)

		if _, err := b.api.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, "internal error"),
		); err != nil {
			b.log.Info("%w", err)
		}
	}
}
