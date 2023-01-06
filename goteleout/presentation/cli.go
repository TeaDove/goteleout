package Presentation

import (
	"github.com/urfave/cli/v2" // imports as package "cli"

	"os"
)

type PresentationCli struct {
}

func (presentatioCli PresentationCli) Run() {
	(&cli.App{}).Run(os.Args)
}
