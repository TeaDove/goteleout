package presentation

import (
	"bufio"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

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

func captureInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			color.White("exiting")
			os.Exit(int(syscall.SIGINT))
		}
	}()
}
