package telegramsupplier

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendErrorsDoNotLeakToken(t *testing.T) {
	t.Parallel()

	const token = "123456:secret-token"

	filePath := filepath.Join(t.TempDir(), "file.txt")

	err := os.WriteFile(filePath, []byte("content"), 0o600)
	require.NoError(t, err)

	testCases := []struct {
		name string
		send func(ctx context.Context, supplier *Supplier) error
	}{
		{
			name: "send message",
			send: func(ctx context.Context, supplier *Supplier) error {
				return supplier.SendMessage(ctx, 1, "text", ModeDefault, false, false)
			},
		},
		{
			name: "send files",
			send: func(ctx context.Context, supplier *Supplier) error {
				return supplier.SendFiles(ctx, 1, []string{filePath}, false)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()

			supplier := NewSupplier(token, new("http://127.0.0.1:1/bot"))

			sendErr := tc.send(tt.Context(), &supplier)

			require.Error(tt, sendErr)
			require.NotContains(tt, sendErr.Error(), token)
			require.Contains(tt, sendErr.Error(), redactedToken)
		})
	}
}
