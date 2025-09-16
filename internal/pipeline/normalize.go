package pipeline

import (
	"github.com/iamlucif3r/trikal/internal/models"
	"github.com/iamlucif3r/trikal/internal/util"
)

func Normalize(in []models.NewsItem) []models.NewsItem {
	for i := range in {
		in[i].Title = util.StripUnsafe(in[i].Title)
		in[i].Summary = util.StripUnsafe(in[i].Summary)
		in[i].Normalize()
	}
	return in
}
