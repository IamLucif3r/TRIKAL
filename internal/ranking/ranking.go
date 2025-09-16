package ranking

import (
	"sort"

	"github.com/iamlucif3r/trikal/internal/config"
	"github.com/iamlucif3r/trikal/internal/models"
)

type Mode int

const (
	ByFinal Mode = iota
	ByReel
	ByDeepDive
	ByExplainer
	ByComposite
)

func Rank(items []models.ScoredNewsItem, mode Mode, topN int, cfg *config.Config) []models.ScoredNewsItem {
	if len(items) == 0 {
		return items
	}

	sort.Slice(items, func(i, j int) bool {
		switch mode {
		case ByReel:
			return items[i].ReelPotential > items[j].ReelPotential
		case ByDeepDive:
			return items[i].DeepDivePotential > items[j].DeepDivePotential
		case ByExplainer:
			return items[i].ExplainerPotential > items[j].ExplainerPotential
		case ByComposite:
			return compositeScore(items[i], cfg) > compositeScore(items[j], cfg)
		default:
			return items[i].FinalScore > items[j].FinalScore
		}
	})

	if topN > 0 && len(items) > topN {
		return items[:topN]
	}
	return items
}

func compositeScore(item models.ScoredNewsItem, cfg *config.Config) float64 {
	return cfg.Ranking.FinalScoreWeight*item.FinalScore +
		cfg.Ranking.TimelinessWeight*float64(item.Timeliness) +
		cfg.Ranking.ReelWeight*float64(item.ReelPotential)
}
