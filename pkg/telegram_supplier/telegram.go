package telegram_supplier

import (
	"fmt"
	"html"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	tele "gopkg.in/telebot.v4"
)

type Supplier struct {
	bot *tele.Bot
}

func NewSupplier(token string) (Supplier, error) {
	bot, err := tele.NewBot(tele.Settings{Token: token})
	if err != nil {
		return Supplier{}, errors.Wrap(err, "create bot")
	}

	return Supplier{bot: bot}, nil
}

func (r *Supplier) SendMessage(
	chatID int64,
	text string,
	parseMode string,
	asCode bool,
	quite bool,
) error {
	if asCode && parseMode != "" {
		color.Yellow(`Settings "parse mode" and "code" simultaneously are not allowed, ignoring parseMode`)
	}

	if asCode {
		parseMode = tele.ModeHTML
		text = fmt.Sprintf("<code>%s</code>", html.EscapeString(text))
	}

	_, err := r.bot.Send(tele.ChatID(chatID), text, &tele.SendOptions{
		DisableNotification: quite,
		ParseMode:           parseMode,
	})
	if err != nil {
		return errors.Wrap(err, "unable to send message")
	}

	return nil
}

func (r *Supplier) SendFiles(chatID int64, filenames []string, quite bool) error {
	for _, filename := range filenames {
		_, err := r.bot.Send(
			tele.ChatID(chatID),
			&tele.Document{File: tele.FromDisk(filename), FileName: filename},
			&tele.SendOptions{
				DisableNotification: quite,
				ParseMode:           tele.ModeHTML,
			})
		if err != nil {
			return errors.Wrap(err, "send file")
		}
	}

	return nil
}
