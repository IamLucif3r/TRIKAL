package pipeline

import (
	"github.com/iamlucif3r/trikal/internal/models"
)

func DedupByURL(in []models.NewsItem) []models.NewsItem {
	seen := make(map[string]struct{}, len(in))
	out := make([]models.NewsItem, 0, len(in))
	for _, it := range in {
		if _, ok := seen[it.URL]; ok {
			continue
		}
		seen[it.URL] = struct{}{}
		out = append(out, it)
	}
	return out
}

func DedupByURLScored(items []models.ScoredNewsItem) []models.ScoredNewsItem {
	seen := make(map[string]struct{}, len(items))
	out := make([]models.ScoredNewsItem, 0, len(items))
	for _, it := range items {
		if _, ok := seen[it.URL]; ok {
			continue
		}
		seen[it.URL] = struct{}{}
		out = append(out, it)
	}
	return out
}
