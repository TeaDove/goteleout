package telegram_supplier

import (
	"fmt"
	"html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type Supplier struct {
	bot *tgbotapi.BotAPI
}

func NewSupplier(token string) (Supplier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return Supplier{}, errors.Wrap(err, "failed to create bot")
	}

	return Supplier{bot: bot}, nil
}

func (r *Supplier) SendMessage(
	chatID int64,
	text string,
	asHTML bool,
	asCode bool,
	quite bool,
) error {
	if !asHTML {
		text = html.EscapeString(text)
	}

	if asCode {
		text = fmt.Sprintf("<code>%s</code>", text)
	}

	message := tgbotapi.NewMessage(chatID, text)
	message.DisableNotification = quite
	message.ParseMode = "html"

	_, err := r.bot.Send(message)
	if err != nil {
		return errors.Wrap(err, "unable to send message")
	}

	return nil
}

func (r *Supplier) SendFiles(chatID int64, filenames []string, quite bool) error {
	for _, filename := range filenames {
		message := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filename))
		message.DisableNotification = quite

		_, err := r.bot.Send(message)
		if err != nil {
			return errors.Wrap(err, "failed to send file")
		}
	}

	return nil
}
