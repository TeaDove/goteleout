package supplier

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teadove/goteleout/internal/schemas"
)

func TestIntegration_BadTokenError(t *testing.T) {
	_, err := NewTelegramSupplier("badToken", true)
	assert.Error(t, err)
}

func TestUnit_GoodTokenOk(t *testing.T) {
	settings, err := schemas.GetSettingsFromFile("")
	assert.NoError(t, err)
	_, err = NewTelegramSupplier(settings.Token, true)
	assert.NoError(t, err)
}

func getSupplierFixture(t *testing.T) (TelegramSupplier, int64) {
	settings, err := schemas.GetSettingsFromFile("")
	assert.NoError(t, err)

	telegramSupplier, err := NewTelegramSupplier(settings.Token, true)
	assert.NoError(t, err)
	userId, err := strconv.ParseInt(settings.User, 10, 64)
	assert.NoError(t, err)
	return telegramSupplier, userId
}

func TestIntegration_SendMessage(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture(t)
	err := telegramSupplier.SendMessage(userId, "test message with <b>tags</b>", false, true, true)
	assert.NoError(t, err)
}

func TestIntegration_SendMessageGoodHTML(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture(t)
	err := telegramSupplier.SendMessage(userId, "test message with <b>tags</b>", true, false, true)
	assert.NoError(t, err)
}

func TestIntegration_SendMessageBadHTML(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture(t)
	err := telegramSupplier.SendMessage(
		userId,
		"test message with <b>tags</b><>",
		true,
		false,
		true,
	)
	assert.Error(t, err)
}

func TestIntegration_SendFiles(t *testing.T) {
	telegramSupplier, userId := getSupplierFixture(t)
	telegramSupplier.SendFiles(userId, "../../go.mod telegram.go", true)
}
