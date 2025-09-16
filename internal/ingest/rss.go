package ingest

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/iamlucif3r/trikal/internal/config"
	"github.com/iamlucif3r/trikal/internal/logging"
	"github.com/iamlucif3r/trikal/internal/models"
)

type rssFeed struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FetchRSSAll(ctx context.Context, httpc *http.Client, cfg *config.Config) ([]models.NewsItem, error) {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		result []models.NewsItem
	)
	logger := logging.New(cfg.Log)
	for _, feed := range cfg.Ingest.RSSFeeds {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				logger.Errorf("rss request build failed for %s: %v", url, err)
				return
			}
			resp, err := httpc.Do(req)
			if err != nil {
				logger.Errorf("rss fetch failed for %s: %v", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Warnf("rss feed %s returned status %d", url, resp.StatusCode)
				return
			}

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Errorf("rss read failed for %s: %v", url, err)
				return
			}

			var rss rssFeed
			if err := xml.Unmarshal(b, &rss); err != nil {
				logger.Errorf("rss parse failed for %s: %v", url, err)
				return
			}

			for _, item := range rss.Channel.Items {
				pub, _ := time.Parse(time.RFC1123Z, item.PubDate)
				ni := models.NewsItem{
					Source:      url,
					SourceType:  models.SourceRSS,
					Title:       strings.TrimSpace(item.Title),
					Summary:     strings.TrimSpace(item.Description),
					URL:         item.Link,
					Authors:     []string{},
					PublishedAt: pub,
					Tags:        []string{},
					Raw:         b,
					Metadata:    map[string]string{},
					ContentHash: "", // fill later in pipeline
					CreatedAt:   time.Now().UTC(),
				}
				mu.Lock()
				result = append(result, ni)
				mu.Unlock()
			}
			logger.Debug("rss fetched %d items from %s", len(rss.Channel.Items), url)
		}(feed)
	}
	wg.Wait()

	logger.Debug("total RSS items collected: %d", len(result))
	return result, nil
}
