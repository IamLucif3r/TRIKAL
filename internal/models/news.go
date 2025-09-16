package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type SourceType string

const (
	SourceRSS SourceType = "rss"
	SourceNVD SourceType = "nvd"
)

type NewsItem struct {
	Source             string            `json:"source"`
	SourceType         SourceType        `json:"source_type"`
	Title              string            `json:"title"`
	Summary            string            `json:"summary"`
	URL                string            `json:"url"`
	Authors            []string          `json:"authors"`
	PublishedAt        time.Time         `json:"published_at"`
	Tags               []string          `json:"tags"`
	Raw                json.RawMessage   `json:"raw"`
	Metadata           map[string]string `json:"metadata"`
	ContentHash        string            `json:"content_hash"`
	CreatedAt          time.Time         `json:"created_at"`
	FinalScore         float64           `json:"final_score"`
	ReelPotential      int               `json:"reel_potential"`
	DeepDivePotential  int               `json:"deep_dive_potential"`
	ExplainerPotential int               `json:"explainer_potential"`
	Timeliness         int               `json:"timeliness"`
	AudienceRelevance  int               `json:"audience_relevance"`
}

func (n *NewsItem) Normalize() {
	n.Title = strings.TrimSpace(n.Title)
	n.Summary = strings.TrimSpace(n.Summary)
	n.URL = strings.TrimSpace(n.URL)
	if n.Metadata == nil {
		n.Metadata = map[string]string{}
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now().UTC()
	}
	key := strings.ToLower(strings.TrimSpace(n.Title)) + "|" + strings.ToLower(strings.TrimSpace(n.URL))
	sum := sha256.Sum256([]byte(key))
	n.ContentHash = hex.EncodeToString(sum[:])
}

func FromRSS(fp *gofeed.Feed, it *gofeed.Item) NewsItem {
	authors := make([]string, 0, len(it.Authors))
	for _, a := range it.Authors {
		if s := strings.TrimSpace(a.Name); s != "" {
			authors = append(authors, s)
		}
	}
	pub := time.Now().UTC()
	if it.PublishedParsed != nil {
		pub = it.PublishedParsed.UTC()
	} else if it.UpdatedParsed != nil {
		pub = it.UpdatedParsed.UTC()
	}
	return NewsItem{
		Source:      strings.TrimSpace(fp.Title),
		SourceType:  SourceRSS,
		Title:       strings.TrimSpace(it.Title),
		Summary:     strings.TrimSpace(prefer(it.Description, it.Content)),
		URL:         strings.TrimSpace(it.Link),
		Authors:     authors,
		PublishedAt: pub,
		Tags:        nil,
		Metadata:    map[string]string{},
	}
}

func prefer(a, b string) string {
	if t := strings.TrimSpace(a); t != "" {
		return t
	}
	return strings.TrimSpace(b)
}
