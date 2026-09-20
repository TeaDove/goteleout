package clipresentation

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
)

const (
	settingsPathFmt = "%s/.config/goteleout.json"
)

type Settings struct {
	Token string  `json:"token"`
	User  int64   `json:"user"`
	Proxy *string `json:"proxy,omitempty"`
}

func getSettings() (Settings, error) {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return Settings{}, errors.WithStack(err)
	}

	path := fmt.Sprintf(settingsPathFmt, dirName)

	var settings Settings

	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return Settings{}, errors.Wrap(err, "read settings file")
	}

	err = json.Unmarshal(jsonFile, &settings)
	if err != nil {
		return Settings{}, errors.Wrap(err, "json unmarshal")
	}

	if settings.Token == "" || settings.User == 0 {
		return Settings{}, errors.New("either token or user is empty")
	}

	return settings, nil
}
