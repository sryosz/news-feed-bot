package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"news-feed-bot/internal/bot"
	"news-feed-bot/internal/bot/middleware"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/config"
	fetcher "news-feed-bot/internal/fetcher"
	"news-feed-bot/internal/notifier"
	"news-feed-bot/internal/storage"
	"news-feed-bot/internal/summary"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Error("failed to create bot: %v", err)
		os.Exit(1)
	}

	articleStorage, err := storage.NewArticleStorage(log)
	if err != nil {
		log.Error("failed to create article storage", err)
		os.Exit(1)
	}

	sourceStorage, err := storage.NewSourceStorage(log)
	if err != nil {
		log.Error("failed to create source storage", err)
		os.Exit(1)
	}

	f := fetcher.New(
		articleStorage,
		sourceStorage,
		cfg.FetchInterval,
		cfg.FilterKeywords,
		log,
	)

	n := notifier.New(
		articleStorage,
		summary.New(cfg.OpenAIKey, cfg.OpenAIPrompt),
		botAPI,
		cfg.NotificationInterval,
		2*cfg.FetchInterval,
		cfg.TelegramChannelID,
		log,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	//separate registering views for bot somehow
	newsBot := botkit.New(log, botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart())
	newsBot.RegisterCmdView("addsource",
		middleware.AdminOnly(
			cfg.TelegramChannelID,
			bot.ViewCmdAddSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView("listsources",
		middleware.AdminOnly(
			cfg.TelegramChannelID,
			bot.ViewCmdListSources(sourceStorage),
		),
	)
	newsBot.RegisterCmdView("commands",
		middleware.AdminOnly(
			cfg.TelegramChannelID,
			bot.ViewCmdListCommands(newsBot.CmdViews),
		),
	)

	go func(ctx context.Context) {
		if err := f.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Error("failed to start fetcher", err)
				return
			}

			log.Info("fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := n.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Error("failed to start notifier", err)
				return
			}

			log.Info("notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		log.Error("failed to run botkit", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
