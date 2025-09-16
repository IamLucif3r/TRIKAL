package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"time"

	"github.com/iamlucif3r/trikal/internal/config"
	"github.com/iamlucif3r/trikal/internal/logging"
	"github.com/iamlucif3r/trikal/internal/models"
)

type nvdResponse struct {
	Vulnerabilities []struct {
		Cve struct {
			ID           string `json:"id"`
			Published    string `json:"published"`
			LastModified string `json:"lastModified"`
			Descriptions []struct {
				Lang  string `json:"lang"`
				Value string `json:"value"`
			} `json:"descriptions"`
			References []struct {
				URL string `json:"url"`
			} `json:"references"`
		} `json:"cve"`
	} `json:"vulnerabilities"`
}

func FetchNVD(ctx context.Context, httpc *http.Client, since, until time.Time, cfg *config.Config) ([]models.NewsItem, error) {
	api := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?pubStartDate=%s&pubEndDate=%s",
		since.Format("2006-01-02T15:04:05.000Z"),
		until.Format("2006-01-02T15:04:05.000Z"),
	)
	logger := logging.New(cfg.Log)

	req, err := http.NewRequestWithContext(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warn("[WARN] NVD returned %d", resp.StatusCode)
		return nil, nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed nvdResponse
	if err := json.Unmarshal(b, &parsed); err != nil {
		return nil, err
	}

	var result []models.NewsItem
	for _, v := range parsed.Vulnerabilities {
		pub, _ := time.Parse(time.RFC3339, v.Cve.Published)
		summary := ""
		if len(v.Cve.Descriptions) > 0 {
			summary = v.Cve.Descriptions[0].Value
		}
		url := ""
		if len(v.Cve.References) > 0 {
			url = v.Cve.References[0].URL
		}

		result = append(result, models.NewsItem{
			Source:      "NVD",
			SourceType:  models.SourceNVD,
			Title:       v.Cve.ID,
			Summary:     summary,
			URL:         url,
			Authors:     []string{},
			PublishedAt: pub,
			Tags:        []string{"CVE"},
			Raw:         b,
			Metadata:    map[string]string{"lastModified": v.Cve.LastModified},
			ContentHash: "",
			CreatedAt:   time.Now().UTC(),
		})
	}

	return result, nil
}
