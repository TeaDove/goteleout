package schemas

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Settings(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "goteleout_test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	err = os.WriteFile(file.Name(), []byte("{\"token\": \"abc\", \"user\": \"123\"}"), 0644)
	assert.NoError(t, err)

	settings, err := GetSettingsFromFile(file.Name())
	assert.NoError(t, err)

	assert.Equal(t, settings.User, "123", "Default user assertion error")
	assert.Equal(t, settings.Token, "abc", "Token assertion error")
}
