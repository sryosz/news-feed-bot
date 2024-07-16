package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
	"strconv"
	"strings"
)

func ViewCmdListCommands(commands map[string]botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		// TODO: add description for each command & make pretty output
		var commandsName []string

		i := 1
		for k := range commands {
			commandsName = append(
				commandsName,
				strconv.Itoa(i)+"\\. /"+k,
			)
			i++
		}

		msgText := fmt.Sprintf(
			"List of commands: \n\n%s",
			strings.Join(commandsName, "\n\n"),
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
