package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/iamlucif3r/trikal/internal/database"
	"github.com/iamlucif3r/trikal/internal/types"
	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
)

func FetchNews() error {
	var rssSources []string

	yamlFile, err := ioutil.ReadFile("rss.yaml")
	if err != nil {
		return fmt.Errorf("error loading 'rss.yaml' file : %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &rssSources)
	if err != nil {
		return fmt.Errorf("error unmarshaling 'rss.yaml' file : %v", err)
	}
	for _, url := range rssSources {
		log.Println("[INFO] Reading ", url, " to fetch news")
		err := fetchAndStoreFromURL(url)
		if err != nil {
			log.Printf("[Error] Error fetching news from %s: %v\n", url, err)
			continue
		}
	}
	log.Println("[SUCCESS] All news fetched successfully from RSS sources")
	return nil
}

func fetchAndStoreFromURL(feedURL string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return fmt.Errorf("failed to parse RSS: %v", err)
	}

	for _, item := range feed.Items {
		news := types.NewsItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PublishedAt: item.PublishedParsed.Format(time.RFC3339),
			Source:      feed.Title,
		}

		if err := insertArticle(news); err != nil {
			log.Printf("[Error] Error inserting news: %v\n", err)
		}
	}

	return nil
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
