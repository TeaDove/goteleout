package presentation

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/teadove/goteleout/pkg/telegram_supplier"
	"github.com/urfave/cli/v3"
	tele "gopkg.in/telebot.v4"
)

func action(_ context.Context, c *cli.Command) error {
	settings, err := getSettings()
	if err != nil {
		innerErr := setDefaultSettings()
		if innerErr != nil {
			return errors.Wrap(innerErr, "set default settings")
		}

		return errors.Wrap(err, "set settings, edit them at ~/.config/teleout.json")
	}

	telegramSupplier, err := telegram_supplier.NewSupplier(settings.Token)
	if err != nil {
		return errors.Wrap(err, "new supplier")
	}

	messageText, err := readFromPipe()
	if err != nil {
		messageText = strings.Join(c.Args().Slice(), " ")
	}

	if messageText == "" {
		messageText = "Hello World!\n\nWith Love from teleout"
	}

	if c.Bool(fileArg) {
		err = telegramSupplier.SendFiles(settings.User, strings.Fields(messageText), c.Bool(quiteArg))
		if err != nil {
			return errors.Wrap(err, "send files")
		}

		return nil
	}

	err = telegramSupplier.SendMessage(
		settings.User,
		messageText,
		c.String(parseModeArg),
		c.Bool(codeArg),
		c.Bool(quiteArg),
	)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return nil
}

const (
	quiteArg     = "quite"
	codeArg      = "code"
	parseModeArg = "parse-mode"
	fileArg      = "file"
)

func Run() {
	captureInterrupt()

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    codeArg,
			Aliases: []string{"c"},
			Value:   false,
			Usage:   "send text with <code> tag to make it monospace, automatically set parseMode=HTML and escapes content",
		},
		&cli.BoolFlag{
			Name:    quiteArg,
			Aliases: []string{"q"},
			Value:   false,
			Usage:   "send message without notifications",
		},
		&cli.StringFlag{
			Name:  parseModeArg,
			Value: tele.ModeDefault,
			Usage: fmt.Sprintf("sets parse mode, can be: %s, %s, %s", tele.ModeHTML, tele.ModeMarkdown, tele.ModeMarkdownV2),
		},
		&cli.BoolFlag{
			Name:    fileArg,
			Aliases: []string{"f"},
			Value:   false,
			Usage:   "specify files to send",
		},
	}

	app := &cli.Command{
		Name:      "goteleout",
		Usage:     "pipes data to telegram, https://github.com/teadove/goteleout",
		UsageText: "goteleout [options] \"text to send\"",
		Flags:     flags,
		Action:    action,
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		color.Red("Unexpected error during execution\n")
		color.White(err.Error())
		os.Exit(1)
	}
}
