package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/iamlucif3r/trikal/internal/models"
)

const cisaURL = "https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json"

type cisaFeed struct {
	Vulnerabilities []struct {
		CveID              string `json:"cveID"`
		VendorProject      string `json:"vendorProject"`
		Product            string `json:"product"`
		VulnerabilityName  string `json:"vulnerabilityName"`
		ShortDescription   string `json:"shortDescription"`
		RequiredAction     string `json:"requiredAction"`
		DueDate            string `json:"dueDate"`
		KnownRansomwareUse string `json:"knownRansomwareCampaignUse"`
		Notes              string `json:"notes"`
		DateAdded          string `json:"dateAdded"`
	} `json:"vulnerabilities"`
}

func FetchCISA(ctx context.Context, httpc *http.Client) ([]models.NewsItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cisaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CISA KEV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from CISA: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var feed cisaFeed
	if err := json.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to decode CISA JSON: %w", err)
	}

	var items []models.NewsItem
	now := time.Now().UTC()
	for _, v := range feed.Vulnerabilities {
		publishedAt, _ := time.Parse("2006-01-02", v.DateAdded)

		metadata := map[string]string{
			"vendor":             v.VendorProject,
			"product":            v.Product,
			"required_action":    v.RequiredAction,
			"due_date":           v.DueDate,
			"known_ransomware":   v.KnownRansomwareUse,
			"notes":              v.Notes,
			"cve_id":             v.CveID,
			"vulnerability_name": v.VulnerabilityName,
		}

		item := models.NewsItem{
			Source:      "CISA KEV",
			SourceType:  models.SourceType("cisa"),
			Title:       fmt.Sprintf("%s - %s", v.CveID, v.VulnerabilityName),
			Summary:     v.ShortDescription,
			URL:         "https://www.cisa.gov/known-exploited-vulnerabilities-catalog",
			Authors:     []string{"CISA"},
			PublishedAt: publishedAt,
			Tags:        []string{"kev", "cisa", "cve"},
			Metadata:    metadata,
			CreatedAt:   now,
		}

		item.ContentHash = ""

		items = append(items, item)
	}

	return items, nil
}
