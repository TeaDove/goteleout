package presentation

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/goteleout/pkg/telegram_supplier"
	"github.com/teadove/teasutils/utils/logger_utils"

	"github.com/urfave/cli/v3"
)

func action(_ context.Context, c *cli.Command) error {
	settings, err := getSettings()
	if err != nil {
		innerErr := setDefaultSettings()
		if innerErr != nil {
			return errors.Wrap(innerErr, "failed to set default settings")
		}

		return errors.Wrap(err, "failed to set settings, edit them at ~/.config/teleout.json")
	}

	telegramSupplier, err := telegram_supplier.NewSupplier(settings.Token)
	if err != nil {
		return errors.WithStack(err)
	}

	messageText, err := readFromPipe()
	if err != nil {
		messageText = strings.Join(c.Args().Slice(), " ")
	}

	if messageText == "" {
		messageText = "Hello World!\n\nWith Love from teleout"
	}

	userID, err := strconv.ParseInt(settings.User, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse user id")
	}

	if c.Bool(fileArg) {
		err = telegramSupplier.SendFiles(userID, strings.Fields(messageText), c.Bool(quiteArg))
		if err != nil {
			return errors.Wrap(err, "failed to send files")
		}

		return nil
	}

	err = telegramSupplier.SendMessage(
		userID,
		messageText,
		c.Bool(htmlArg),
		c.Bool(codeArg),
		c.Bool(quiteArg),
	)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}

const quiteArg = "quite"
const codeArg = "code"
const htmlArg = "html"
const fileArg = "file"

func Run() {
	captureInterrupt()

	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    codeArg,
			Aliases: []string{"c"},
			Value:   false,
			Usage:   "send text with <code> tag to make it monospace",
		},
		&cli.BoolFlag{
			Name:    quiteArg,
			Aliases: []string{"q"},
			Value:   false,
			Usage:   "send message without notifications",
		},
		&cli.BoolFlag{
			Name:  htmlArg,
			Value: false,
			Usage: "do no escape html tags",
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

	ctx := logger_utils.NewLoggedCtx()
	err := app.Run(ctx, os.Args)

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error during execution")

		os.Exit(1)
	}
}
