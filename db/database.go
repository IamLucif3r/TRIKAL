package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/IamLucif3r/trikal/models"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("[Error] Cannot connect to DB: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("[Error] Cannot ping DB: %v", err)
	}
	log.Println("[Info] Connected to DB")
}

func InsertArticle(ctx context.Context, article models.Article) error {
	query := `
	INSERT INTO articles (title, link, summary_raw, source, published_at)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (link) DO NOTHING;
	`

	_, err := DB.ExecContext(ctx, query,
		article.Title,
		article.Link,
		article.SummaryRaw,
		article.Source,
		article.PublishedAt,
	)

	return err
}
