package ports

import (
	"context"

	"ehgm.com.br/url-shortener/domain/model"
)

type UrlService interface {
	GenerateId(ctx context.Context, url string) (string, error)
	GetUrl(ctx context.Context, id string) (*model.ShortUrl, error)
	GetUrlToRedirect(ctx context.Context, id string) (string, bool, error)
	UpdateUrl(ctx context.Context, id string, json map[string]interface{}) error
	GetStats(ctx context.Context, limit int) ([]model.ShortUrl, error)
}
