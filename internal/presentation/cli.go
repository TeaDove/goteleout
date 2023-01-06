package presentation

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

type PresentationCli struct {
}

func (presentatioCli PresentationCli) Run() {
	var err = (&cli.App{}).Run(os.Args)
	if err != nil {
		fmt.Printf("Error occured: %s", err)
	}
}
