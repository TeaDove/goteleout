package shared

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const defaultPath = "%s/.config/teleout.json"

type ClientSettings struct {
	Token       string `json:"token"`
	DefaultUser string `json:"user"`
}

func GetSettingsFromFile(path string) (ClientSettings, error) {
	if path == "" {
		dirName, err := os.UserHomeDir()
		if err != nil {
			return ClientSettings{}, err
		}
		path = fmt.Sprintf(defaultPath, dirName)
	}
	var settings ClientSettings

	jsonFile, err := os.Open(path)
	if err != nil {
		return ClientSettings{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		return ClientSettings{}, err
	}

	return settings, nil
}
