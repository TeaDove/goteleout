package presentation

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	settingsPathFmt = "%s/.config/teleout.json"
)

type Settings struct {
	Token string `json:"token"`
	User  string `json:"user"`
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
		return Settings{}, errors.Wrapf(err, "failed to read settings file")
	}

	err = json.Unmarshal(jsonFile, &settings)
	if err != nil {
		return Settings{}, errors.WithStack(err)
	}

	if settings.Token == "" || settings.User == "" {
		return Settings{}, errors.New("either token or user is empty")
	}

	return settings, nil
}

func setDefaultSettings() error {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return errors.WithStack(err)
	}

	path := fmt.Sprintf(settingsPathFmt, dirName)

	var settings = Settings{
		Token: "<telegram-token>",
		User:  "<user-id>",
	}

	bytes, err := json.Marshal(settings)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.WriteFile(path, bytes, 0600)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
