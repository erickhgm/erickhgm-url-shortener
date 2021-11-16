package ports

import (
	"context"

	"ehgm.com.br/url-shortener/domain/model"
)

type UrlRepository interface {
	Save(ctx context.Context, id, url string, enable bool) error
	FindById(ctx context.Context, id string) (*model.ShortUrl, error)
	Update(ctx context.Context, id string, json map[string]interface{}) error
	GetStats(ctx context.Context, limit int) ([]model.ShortUrl, error)
}
