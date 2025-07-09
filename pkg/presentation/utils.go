package presentation

import (
	"bufio"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/teadove/teasutils/utils/logger_utils"
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
		return "", errors.Wrap(err, "unable to copy from buf")
	}

	return buf.String(), nil
}

func captureInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			zerolog.Ctx(logger_utils.NewLoggedCtx()).
				Info().
				Stringer("signal", sig).
				Msg("exiting")

			os.Exit(int(syscall.SIGINT))
		}
	}()
}
