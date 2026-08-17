package telegramsupplier

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/cockroachdb/errors"
	"github.com/fatih/color"
)

const (
	sendMessageTimeout = 500 * time.Millisecond
	sendFilesTimeout   = 10 * time.Second
)

func (r *Supplier) SendMessage(
	ctx context.Context,
	chatID int64,
	text, parseMode string,
	asCode, quite bool,
) error {
	if asCode && parseMode != "" {
		color.Yellow(`Settings "parse mode" and "code" simultaneously are not allowed, ignoring parseMode`)
	}

	if asCode {
		parseMode = ModeHTML
		text = fmt.Sprintf("<code>%s</code>", html.EscapeString(text))
	}

	form := url.Values{}
	form.Set("chat_id", strconv.FormatInt(chatID, 10))
	form.Set("text", text)
	if parseMode != "" {
		form.Set("parse_mode", parseMode)
	}
	if quite {
		form.Set("disable_notification", "true")
	}

	err := retry.Do(func() error {
		return r.callForm(ctx, sendMessageTimeout, "sendMessage", form)
	}, retry.Attempts(retryAttempts))
	if err != nil {
		return errors.Wrap(err, "unable to send message")
	}

	return nil
}

func (r *Supplier) SendFiles(ctx context.Context, chatID int64, filenames []string, quite bool) error {
	for _, filename := range filenames {
		err := retry.Do(func() error {
			return r.sendDocument(ctx, chatID, filename, quite)
		}, retry.Attempts(retryAttempts))
		if err != nil {
			return errors.Wrap(err, "send file")
		}
	}

	return nil
}

func (r *Supplier) sendDocument(ctx context.Context, chatID int64, filename string, quite bool) error {
	ctx, cancel := context.WithTimeout(ctx, sendFilesTimeout)
	defer cancel()

	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fields := map[string]string{
		"chat_id":              strconv.FormatInt(chatID, 10),
		"parse_mode":           ModeHTML,
		"disable_notification": strconv.FormatBool(quite),
	}
	for key, value := range fields {
		err = writer.WriteField(key, value)
		if err != nil {
			return errors.Wrap(err, "write field")
		}
	}

	part, err := writer.CreateFormFile("document", filepath.Base(filename))
	if err != nil {
		return errors.Wrap(err, "create form file")
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return errors.Wrap(err, "copy file")
	}

	err = writer.Close()
	if err != nil {
		return errors.Wrap(err, "close writer")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.methodURL("sendDocument"), body)
	if err != nil {
		return errors.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := r.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}
