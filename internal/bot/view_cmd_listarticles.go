package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/botkit/markup"
	"news-feed-bot/internal/model"
	"strings"
)

type ArticleLister interface {
	ArticlesBySourceID(ctx context.Context, sourceID int64) ([]model.Article, error)
}

func ViewCmdListArticles(lister ArticleLister) botkit.ViewFunc {
	const op = "bot.ViewCmdListArticles"

	type listArticlesArgs struct {
		ID int64 `json:"id"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[listArticlesArgs](update.Message.CommandArguments())
		if err != nil {
			log.Println("1")
			return fmt.Errorf("%s: %w", op, err)
		}

		//TODO: the message is too large so split up for some pages (?)
		articles, err := lister.ArticlesBySourceID(ctx, args.ID)
		if err != nil {
			log.Println("2")
			return fmt.Errorf("%s: %w", op, err)
		}

		var articleInfos []string

		for _, article := range articles {
			articleInfos = append(articleInfos, formatArticle(article))
		}

		msgText := fmt.Sprintf(
			"List of articles \\(total %d\\):\n\n%s",
			len(articles),
			strings.Join(articleInfos, "\n\n"),
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			log.Println("3")
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	}
}

func formatArticle(article model.Article) string {
	return fmt.Sprintf(
		"⚪ *[%s](%s)",
		markup.EscapeForMarkdown(article.Title),
		markup.EscapeForMarkdown(article.Link),
	)
}
