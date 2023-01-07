package supplier

import (
	"strconv"
	"testing"

	"github.com/teadove/goteleout/internal/schemas"
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

func getSupplierFixture() (TelegramSupplier, int64) {
	settings, err := schemas.GetSettingsFromFile("")
	if err != nil {
		panic(err)
	}

	telegramSupplier, err := NewTelegramSupplier(settings.Token, true)
	if err != nil {
		panic(err)
	}
	userId, err := strconv.ParseInt(settings.User, 10, 64)
	if err != nil {
		panic(err)
	}
	return telegramSupplier, userId
}

func TestSupplier_SendMessage(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture()
	err := telegramSupplier.SendMessage(userId, "test message with <b>tags</b>", false, true, true)
	if err != nil {
		panic(err)
	}
}

func TestSupplier_SendMessageGoodHTML(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture()
	err := telegramSupplier.SendMessage(userId, "test message with <b>tags</b>", true, false, true)
	if err != nil {
		panic(err)
	}
}

func TestSupplier_SendMessageBadHTML(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture()
	err := telegramSupplier.SendMessage(userId, "test message with <b>tags</b><>", true, false, true)
	if err == nil {
		panic("Should return error!")
	}
}

func TestSupplier_SendFiles(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture()
	telegramSupplier.SendFiles(userId, "../../go.mod telegram.go", true)
}
