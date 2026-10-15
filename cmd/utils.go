package cmd

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

func getCommand(c *cli.Command) string {
	command, _ := readFromPipe()

	command = strings.TrimSpace(command)
	if command == "" {
		command = strings.TrimSpace(strings.Join(c.Args().Slice(), " "))
	}

	if command == "" {
		return "Hello World!"
	}

	return command
}

func readFromPipe() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", errors.New("no data in pipeline")
	}

	reader := bufio.NewReader(os.Stdin)
	buf := new(strings.Builder)

	_, err := io.Copy(buf, reader)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return buf.String(), nil
}

const commandTimeout = 90 * time.Second

func captureInterrupt() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			cancel()
			time.Sleep(50 * time.Millisecond)
			color.White("exiting")
			os.Exit(int(syscall.SIGINT))
		}
	}()

	return ctx
}
