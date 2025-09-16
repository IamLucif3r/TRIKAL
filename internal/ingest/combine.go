package ingest

import (
	"context"
	"net/http"
	"time"

	"github.com/iamlucif3r/trikal/internal/config"
	"github.com/iamlucif3r/trikal/internal/httpx"
	"github.com/iamlucif3r/trikal/internal/models"
)

func NewHTTPClient(cfg *config.Config) *http.Client {
	return httpx.New(cfg.Timeout(), cfg.HTTP.UserAgent, cfg.HTTP.MaxIdlePerHost).Client
}

func FetchAll(ctx context.Context, httpc *http.Client, cfg *config.Config) ([]models.NewsItem, error) {
	rss, err := FetchRSSAll(ctx, httpc, cfg)
	if err != nil {
		return nil, err
	}

	until := time.Now().UTC()
	since := until.Add(-24 * time.Hour)

	nvd, err := FetchNVD(ctx, httpc, since, until, cfg)
	if err != nil {
		return nil, err
	}

	cisa, err := FetchCISA(ctx, httpc)
	if err != nil {
		return nil, err
	}

	all := append(rss, nvd...)
	all = append(all, cisa...)

	return all, nil
}
