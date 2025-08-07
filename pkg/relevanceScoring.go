package pkg

import (
	"fmt"
	"log"
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
	You are a cybersecurity educator and strategist building in-depth, high-value content for the brand "pwnspectrum".

	You publish:
	- YouTube explainers titled "How X Got Hacked"
	- CVE walkthroughs with PoC, root cause, and exploitation details
	- Medium articles that break down complex attacks or tools
	- Educational reels summarizing attack chains and tricks

	Given the news article:

	Title: %s
	Description: %s

	Rate how suitable this article is to be turned into **a hands-on walkthrough, technical breakdown, or explanatory content**.
	🔍 Ask yourself:
	- Is this breaking news, or highly relatable for practitioners?
	- Could it be a fun, shocking, or meme-worthy Reel?
	- Could it be explored in depth in a Medium article with research and value?
	- Could it be turned into an engaging YouTube story, walkthrough, or digest segment?

	Prioritize:
	✅ Real-world incidents (breaches, new techniques, active campaigns)
	✅ New vulnerabilities or CVEs (especially with PoCs or attack details)
	✅ Multi-stage attacks or phishing tactics that can be explained visually
	✅ Anything that shows *how* something works or broke

	❌ Ignore purely corporate press releases, product announcements, or vague threat alerts with no technical meat.

	Score from 0 (useless for content) to 10 (perfect for technical breakdown or teaching).

	💡 Imagine you're selecting top content out of 100 headlines this week.
	Only **10 headlines should score 9 or 10.**
	Score high ONLY when it's truly 🔥 — unique, timely, and rich in content.

	🎯 Score high ONLY if it can become a **great "How X Got Hacked" video, CVE walkthrough, or research-backed article**.

	Return ONLY a single number, no text, no comments, no formatting.
	`, news.Title, news.Description)

	scoreStr, err := QueryOllamaAPI(prompt)
	if err != nil {
		log.Println("Error Querying Ollama API: ", err)
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
