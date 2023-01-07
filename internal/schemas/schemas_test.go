package schemas

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShared_Settings(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "goteleout_test")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	err = os.WriteFile(file.Name(), []byte("{\"token\": \"abc\", \"user\": \"123\"}"), 0644)
	if err != nil {
		panic(err)
	}

	settings, err := GetSettingsFromFile(file.Name())
	if err != nil {
		panic(err)
	}
	t.Logf("%+v\n", settings)

	assert.Equal(t, settings.User, "123", "Default user assertion error")
	assert.Equal(t, settings.Token, "abc", "Token assertion error")
}
