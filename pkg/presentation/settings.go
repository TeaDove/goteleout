package presentation

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	settingsPathFmt = "%s/.config/goteleout.json"
)

type Settings struct {
	Token string `json:"token"`
	User  int64  `json:"user"`
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
		return Settings{}, errors.Wrapf(err, "read settings file")
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

func setDefaultSettings() error {
	dirName, err := os.UserHomeDir()
	if err != nil {
		return errors.WithStack(err)
	}

	path := fmt.Sprintf(settingsPathFmt, dirName)

	var settings = Settings{
		Token: "TelegramBotToken",
		User:  123,
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
