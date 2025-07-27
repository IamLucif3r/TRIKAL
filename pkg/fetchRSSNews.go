package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/iamlucif3r/trikal/internal/database"
	"github.com/iamlucif3r/trikal/internal/types"
	"github.com/mmcdole/gofeed"
)

var cybersecurityKeywords = []string{
	"exploit", "zero day", "vulnerability", "rce", "breach", "malware",
	"phishing", "ransomware", "APT", "Critical", "CVE", "patch", "cyber attack",
	"rootkit", "leak", "ddos", "security advisory", "mitre", "cisa", "cert",
}

func FetchNews() error {
	// var rssSources []string

	// yamlFile, err := ioutil.ReadFile("rss.yaml")
	// if err != nil {
	// 	return fmt.Errorf("error loading 'rss.yaml' file : %v", err)
	// }

	// err = yaml.Unmarshal(yamlFile, &rssSources)
	// if err != nil {
	// 	return fmt.Errorf("error unmarshaling 'rss.yaml' file : %v", err)
	// }
	// var count int
	// for _, url := range rssSources {
	// 	log.Println("[INFO] Reading ", url, " to fetch news")
	// 	count, err = fetchAndStoreFromURL(url)
	// 	if err != nil {
	// 		log.Printf("[Error] Error fetching news from %s: %v\n", url, err)
	// 		continue
	// 	}
	// }
	// log.Println("[SUCCESS] Fetched [", count, "]  articles successfully from RSS sources")

	// err = ScoreAllArticlesWithLLM()
	// if err != nil {
	// 	log.Printf("[Error] Error scoring articles with LLM: %v\n", err)
	// }
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	err := SendAlertToDiscord(webhookURL)
	if err != nil {
		log.Printf("[Error] Error sending top articles to Discord: %v\n", err)
	}

	return nil
}

func fetchAndStoreFromURL(feedURL string) (int, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse RSS: %v", err)
	}
	var article int
	for _, item := range feed.Items {

		if !passesKeywordFilter(item.Title, item.Description) {
			continue
		}

		news := types.NewsItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PublishedAt: item.PublishedParsed.Format(time.RFC3339),
			Source:      feed.Title,
		}
		article++
		if err := insertArticle(news); err != nil {
			log.Printf("[Error] Error inserting news: %v\n", err)
			article--
		}

	}

	return article, nil
}

func insertArticle(news types.NewsItem) error {
	db := database.DB
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	query := `
	INSERT INTO articles (title, description, link, published_at, source)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (link) DO NOTHING;
	`

	_, err := db.Exec(query,
		news.Title,
		news.Description,
		news.Link,
		news.PublishedAt,
		news.Source,
	)

	return err
}

func passesKeywordFilter(title, description string) bool {
	content := strings.ToLower(title + " " + description)
	for _, keyword := range cybersecurityKeywords {
		if strings.Contains(content, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func ScoreAllArticlesWithLLM() error {
	log.Println("[INFO] Scoring all articles with LLM...")
	articles, err := GetAllArticles()
	if err != nil {
		return err
	}
	for _, article := range articles {
		score, err := ScoreWithLLM(article)
		if err != nil {
			log.Printf("Failed scoring article %d: %v", article.ID, err)
			continue
		}
		err = UpdateLLMScore(article.ID, score)
		if err != nil {
			log.Printf("Failed updating LLM score for article %d: %v", article.ID, err)
		}
	}
	return nil
}

func GetTopArticles(limit int) ([]types.NewsItem, error) {
	db := database.DB
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	rows, err := db.Query(
		"SELECT id, title, description, link, published_at, source, llm_score FROM articles ORDER BY llm_score DESC, published_at DESC LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.NewsItem
	for rows.Next() {
		var news types.NewsItem
		if err := rows.Scan(&news.ID, &news.Title, &news.Description, &news.Link, &news.PublishedAt, &news.Source, &news.LLMScore); err != nil {
			return nil, err
		}
		items = append(items, news)
	}
	return items, nil
}
