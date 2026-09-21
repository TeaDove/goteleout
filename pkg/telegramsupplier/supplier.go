package telegramsupplier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	ModeDefault    = ""
	ModeHTML       = "HTML"
	ModeMarkdown   = "Markdown"
	ModeMarkdownV2 = "MarkdownV2"
)

const (
	apiBaseURL    = "https://api.telegram.org/bot"
	retryAttempts = 5
	redactedToken = "[REDACTED]"
)

type Supplier struct {
	token   string
	baseURL string
	client  *http.Client
}

func NewSupplier(token string, proxy *string) Supplier {
	baseURL := apiBaseURL
	if proxy != nil {
		baseURL = *proxy
	}

	return Supplier{token: token, baseURL: baseURL, client: &http.Client{}}
}

func (r *Supplier) methodURL(method string) string {
	return fmt.Sprintf("%s%s/%s", r.baseURL, r.token, method)
}

func (r *Supplier) redactToken(err error) error {
	if r.token == "" {
		return err
	}

	return errors.New(strings.ReplaceAll(err.Error(), r.token, redactedToken))
}

func (r *Supplier) callForm(ctx context.Context, timeout time.Duration, method string, form url.Values) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		r.methodURL(method),
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return errors.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}

func checkResponse(resp *http.Response) error {
	var body struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"` //nolint: tagliatelle // as expected
		Description string `json:"description"`
	}

	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return errors.Wrapf(err, "decode response, status %d", resp.StatusCode)
	}

	if !body.OK {
		return errors.Errorf("telegram api error %d: %s", body.ErrorCode, body.Description)
	}

	return nil
}
