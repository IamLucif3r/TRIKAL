package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/IamLucif3r/trikal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type ArticleMeta struct {
	ID    uuid.UUID
	Link  string
	Title string
}

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
func GetArticlesWithoutFullText(ctx context.Context) ([]ArticleMeta, error) {
	query := `
	SELECT id, link, title FROM articles
	WHERE full_text IS NULL OR full_text = ''
	ORDER BY fetched_at DESC
	LIMIT 10;
	`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ArticleMeta
	for rows.Next() {
		var a ArticleMeta
		if err := rows.Scan(&a.ID, &a.Link, &a.Title); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	return results, nil
}

func UpdateFullText(ctx context.Context, id uuid.UUID, text string) error {
	query := `UPDATE articles SET full_text = $1 WHERE id = $2;`
	_, err := DB.ExecContext(ctx, query, text, id)
	return err
}
