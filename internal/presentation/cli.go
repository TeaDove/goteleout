package presentation

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"

	"goteleout/internal/schemas"
	"goteleout/internal/supplier"

	"github.com/urfave/cli/v2"
)

const errorMsg = `Error while accuring settings:

%s\n

You should put json in format
{"token": <telegram-token>, "user": <telegram-user-id>}
in settings file(~/.config/teleout.json), or use cli args

`

func readFromPipe() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", errors.New("no data in pipeline")
	}

	reader := bufio.NewReader(os.Stdin)
	buf := new(strings.Builder)

	_, err := io.Copy(buf, reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getSettings(cCtx *cli.Context) (schemas.ClientSettings, error) {
	user := cCtx.String(userArg)
	token := cCtx.String(tokenArg)

	if user == "" || token == "" {
		settinsPath := cCtx.String(settingsFileArg)
		settings, err := schemas.GetSettingsFromFile(settinsPath)
		return settings, err
	}
	return schemas.ClientSettings{User: user, Token: token}, nil
}

func action(cCtx *cli.Context) error {
	settings, err := getSettings(cCtx)
	if settings.Token == "" {
		fmt.Printf(errorMsg, err)
		os.Exit(1)
	}
	telegramSupplier, err := supplier.NewTelegramSupplier(settings.Token, cCtx.Bool(verboseArg))
	if err != nil {
		panic(err)
	}
	if cCtx.Bool(getMeArg) {
		err := telegramSupplier.GetMe(cCtx.Bool(quiteArg))
		return err
	}

	if settings.User == "" {
		fmt.Printf(errorMsg, err)
		os.Exit(1)
	}
	messageText, err := readFromPipe()
	if err != nil {
		messageText = strings.Join(cCtx.Args().Slice(), " ")
	}
	if messageText == "" {
		messageText = "Hello World!\n\nWith Love from teleout"
	}

	userId, err := strconv.ParseInt(settings.User, 10, 64)
	if err != nil {
		panic(err)
	}
	err = telegramSupplier.SendMessage(userId, messageText, cCtx.Bool(htmlArg), cCtx.Bool(codeArg), cCtx.Bool(quiteArg))

	return err
}

const quiteArg = "quite"
const codeArg = "code"
const htmlArg = "html"
const tokenArg = "token"
const userArg = "user"
const settingsFileArg = "settings-file"
const verboseArg = "verbose"
const getMeArg = "get-me"

func captureInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
}

func RunCli() {
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
		&cli.StringFlag{
			Name:    tokenArg,
			EnvVars: []string{"TELEGRAM_TOKEN"},
			Value:   "",
			Usage:   "telegram api token",
		},
		&cli.StringFlag{
			Name:    userArg,
			Aliases: []string{"u"},
			Value:   "",
			Usage:   "telegram user id",
		},
		&cli.StringFlag{
			Name:  settingsFileArg,
			Value: "",
			Usage: "file to store default settings",
		},
		&cli.BoolFlag{
			Name:    verboseArg,
			Aliases: []string{"v"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:  getMeArg,
			Usage: "will listen for updates, and reply with user_id and chat_id",
			Value: false,
		},
	}

	app := &cli.App{
		Name:      "goteleout",
		Usage:     "pipes data to telegram, https://github.com/teadove/goteleout",
		UsageText: "goteleout [options] \"text to send\"",
		Flags:     flags,
		Action:    action,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

}
