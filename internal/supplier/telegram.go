package supplier

import (
	"fmt"
	"html"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSupplier struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramSupplier(token string, debug bool) (TelegramSupplier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return TelegramSupplier{}, err
	}
	bot.Debug = debug
	return TelegramSupplier{bot: bot}, nil
}

func (telegramSupplier TelegramSupplier) SendMessage(chat_id int64, text string, as_html bool, as_code bool, quite bool) error {
	if !as_html {
		text = html.EscapeString(text)
	}
	if as_code {
		text = fmt.Sprintf("<code>%s</code>", text)
	}

	message := tgbotapi.NewMessage(chat_id, text)
	message.BaseChat.DisableNotification = quite
	message.ParseMode = "html"
	_, err := telegramSupplier.bot.Send(message)
	if err != nil {
		return err
	}
	return nil
}

const replyTemplate = `Get me
from user: (@%s, %d)
in chat: (%s, %d)
with message: %s`
const replyTemplateTaged = `Get me
from user: (@%s, <code>%d</code>)
in chat: (%s, <code>%d</code>)
with message: <code>%s</code>`

func compileReply(update tgbotapi.Update, taged bool) string {
	var template string
	if taged {
		template = replyTemplateTaged
	} else {
		template = replyTemplate
	}
	username := update.Message.From.UserName
	if username == "" {
		username = "None"
	}
	chatTitle := update.Message.Chat.Title
	if chatTitle == "" {
		chatTitle = "None"
	}

	return fmt.Sprintf(template,
		username,
		update.Message.From.ID,
		chatTitle,
		update.Message.Chat.ID,
		html.EscapeString(update.Message.Text))
}

func (telegramSupplier TelegramSupplier) GetMe(quite bool) error {
	bot := telegramSupplier.bot
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Print(compileReply(update, false))
			if quite {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, compileReply(update, true))
			msg.ParseMode = "html"
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
