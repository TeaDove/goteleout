package clipresentation

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/fatih/color"
	"github.com/teadove/goteleout/pkg/telegramsupplier"
	"github.com/urfave/cli/v3"
)

func action(ctx context.Context, c *cli.Command) error {
	settings, err := getSettings()
	if err != nil {
		innerErr := setDefaultSettings()
		if innerErr != nil {
			return errors.Wrap(innerErr, "set default settings")
		}

		return errors.Wrap(err, "set settings, edit them at ~/.config/teleout.json")
	}

	telegramSupplier := telegramsupplier.NewSupplier(settings.Token)

	command := getCommand(c)
	if c.Bool(fileArg) {
		err = telegramSupplier.SendFiles(ctx, settings.User, strings.Fields(command), c.Bool(quiteArg))
		if err != nil {
			return errors.Wrap(err, "send files")
		}

		return nil
	}

	err = telegramSupplier.SendMessage(
		ctx,
		settings.User,
		command,
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
	ctx := captureInterrupt()

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
			Value: telegramsupplier.ModeDefault,
			Usage: fmt.Sprintf(
				"sets parse mode, can be: %s, %s, %s",
				telegramsupplier.ModeHTML,
				telegramsupplier.ModeMarkdown,
				telegramsupplier.ModeMarkdownV2,
			),
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

	err := app.Run(ctx, os.Args)
	if err != nil {
		color.Red("Unexpected error during execution\n")
		color.White(err.Error())
		os.Exit(1)
	}
}
