package goteleoutservice

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/teadove/goteleout/pkg/telegramsupplier"
)

type Service struct {
	user       int64
	tgSupplier telegramsupplier.Supplier
}

func NewService() (*Service, error) {
	cfg, err := getCfg()
	if err != nil {
		return nil, errors.Wrap(err, "get cfg, edit it at ~/.config/goteleout.json")
	}

	tgSupplier := telegramsupplier.NewSupplier(cfg.Token, cfg.Proxy)

	return &Service{user: cfg.User, tgSupplier: tgSupplier}, nil
}

func (r *Service) SendMessage(ctx context.Context, text, parseMode string, asCode, quite bool) error {
	return r.tgSupplier.SendMessage(ctx, r.user, text, parseMode, asCode, quite)
}

func (r *Service) SendFiles(ctx context.Context, paths []string, quite bool) error {
	return r.tgSupplier.SendFiles(ctx, r.user, paths, quite)
}
