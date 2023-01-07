package supplier

import (
	"goteleout/internal/schemas"
	"testing"
)

func TestSupplier_BadTokenError(t *testing.T) {
	_, err := NewTelegramSupplier("badToken", true)
	if err == nil {
		panic("NewTelegramSupplier should return error!")
	}
}

func TestSupplier_GoodTokenOk(t *testing.T) {
	settings, err := schemas.GetSettingsFromFile("")
	if err != nil {
		panic(err)
	}
	_, err = NewTelegramSupplier(settings.Token, true)
	if err != nil {
		panic(err)
	}
}
