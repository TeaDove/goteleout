package goteleoutservice

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
)

const (
	cfgPathFmt = "%s/.config/goteleout.json"
)

type Config struct {
	Token string  `json:"token"`
	User  int64   `json:"user"`
	Proxy *string `json:"proxy,omitempty"`
}

func getCfg() (Config, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return Config{}, errors.WithStack(err)
	}

	path := fmt.Sprintf(cfgPathFmt, dirName)

	var cfg Config

	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return Config{}, errors.Wrap(err, "read cfg file")
	}

	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "json unmarshal")
	}

	if cfg.Token == "" || cfg.User == 0 {
		return Config{}, errors.New("either token or user is empty")
	}

	return cfg, nil
}
