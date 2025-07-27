package pkg

import (
	"fmt"
	"math"

	"github.com/iamlucif3r/trikal/internal/database"
	"github.com/iamlucif3r/trikal/internal/types"
)

func GetAllArticles() ([]types.NewsItem, error) {
	db := database.DB
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	rows, err := db.Query("SELECT id, title, description, link, published_at, source FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.NewsItem
	for rows.Next() {
		var news types.NewsItem
		if err := rows.Scan(&news.ID, &news.Title, &news.Description, &news.Link, &news.PublishedAt, &news.Source); err != nil {
			return nil, err
		}
		items = append(items, news)
	}
	return items, nil
}

func ScoreWithLLM(news types.NewsItem) (int, error) {
	prompt := fmt.Sprintf(`
	You are an expert cybersecurity content strategist, who creates for the following formats: 
	a) Short, viral Reels or Memes for Instagram
	b) In-depth and informative articles for Medium
	c) Engaging, visual YouTube Videos (either weekly cybersecurity news digests, walkthroughs of new tools, or step-by-step deep dives into vulnerabilities/CVEs)

	Given the following news article:

	Title: %s
	Description: %s

	Rate how well this article can be turned into HIGH-VALUE content for those platforms.
	Consider:

	- Would this make a great, timely, or viral REEL or meme (shocking, funny, quick insight, shareable)?
	- Does it offer enough depth, context, or new info for a full Medium ARTICLE?
	- Is it suitable for a YOUTUBE VIDEO -- either as part of a weekly news roundup, a hands-on walkthrough for a new tool/technique, or a detailed exploration of a trending vulnerability or CVE?

	Score from 0 (not useful for any format as content) to 10 (amazing: perfect for Reels, Medium articles, and/or engaging videos).

	Output ONLY a single integer or decimal number (e.g., 8 or 8.5), nothing else. Do NOT say "out of 10" or add any label.
`, news.Title, news.Description)

	scoreStr, err := QueryOllamaAPI(prompt)
	if err != nil {
		return 0, err
	}
	scoreRoundoff := int(math.Ceil(float64(scoreStr)))

	return scoreRoundoff, nil
}

func UpdateLLMScore(id int, score int) error {
	db := database.DB
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	_, err := db.Exec("UPDATE articles SET llm_score=$1 WHERE id=$2", score, id)
	return err
}
