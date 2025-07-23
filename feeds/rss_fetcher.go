package feeds

import (
	"context"
	"log"

	"github.com/IamLucif3r/trikal/db"
	"github.com/IamLucif3r/trikal/models"

	"github.com/mmcdole/gofeed"
)

var RSSFeeds = []string{
	"https://thehackernews.com/feeds/posts/default",
	"https://www.bleepingcomputer.com/feed/",
	"https://www.darkreading.com/rss.xml",
}

func FetchAndStoreArticles() {
	parser := gofeed.NewParser()

	for _, feedURL := range RSSFeeds {
		feed, err := parser.ParseURL(feedURL)
		if err != nil {
			log.Printf("[Error] Failed to parse feed: %s - %v\n", feedURL, err)
			continue
		}

		for _, item := range feed.Items {
			article := models.Article{
				Title:       item.Title,
				Link:        item.Link,
				SummaryRaw:  item.Description,
				Source:      feed.Title,
				PublishedAt: item.PublishedParsed,
			}

			err := db.InsertArticle(context.Background(), article)
			if err != nil {
				log.Printf("[Debug] Skipping duplicate or failed insert: %s\n", item.Link)
			} else {
				log.Printf("[Info] Stored: %s\n", item.Title)
			}
		}
	}
}
