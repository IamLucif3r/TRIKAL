package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/iamlucif3r/trikal/internal/models"
	"github.com/lib/pq"
)

type NewsRepo struct {
	db *sql.DB
}

func NewNewsRepo(db *sql.DB) *NewsRepo {
	return &NewsRepo{db: db}
}

func (r *NewsRepo) UpsertBatch(ctx context.Context, items []models.ScoredNewsItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
	INSERT INTO news_items 
	(source, source_type, title, summary, url, authors, published_at, tags, metadata, content_hash, final_score, reel_potential, deep_dive_potential, explainer_potential, timeliness, audience_relevance, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	ON CONFLICT (content_hash) DO UPDATE SET
		source = EXCLUDED.source,
		source_type = EXCLUDED.source_type,
		title = EXCLUDED.title,
		summary = EXCLUDED.summary,
		url = EXCLUDED.url,
		authors = EXCLUDED.authors,
		published_at = EXCLUDED.published_at,
		tags = EXCLUDED.tags,
		metadata = EXCLUDED.metadata,
		final_score = EXCLUDED.final_score,
		reel_potential = EXCLUDED.reel_potential,
		deep_dive_potential = EXCLUDED.deep_dive_potential,
		explainer_potential = EXCLUDED.explainer_potential,
		timeliness = EXCLUDED.timeliness,
		audience_relevance = EXCLUDED.audience_relevance,
		created_at = EXCLUDED.created_at;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {

		authors := item.Authors
		if authors == nil {
			authors = []string{}
		}
		tags := item.Tags
		if tags == nil {
			tags = []string{}
		}

		metadataJSON, err := json.Marshal(item.Metadata)
		if err != nil || string(metadataJSON) == "null" {
			metadataJSON = []byte(`{}`)
		}

		_, err = stmt.ExecContext(
			ctx,
			item.Source,
			item.SourceType,
			item.Title,
			item.Summary,
			item.URL,
			pq.Array(authors),
			item.PublishedAt,
			pq.Array(tags),
			string(metadataJSON),
			item.ContentHash,
			item.FinalScore,
			item.ReelPotential,
			item.DeepDivePotential,
			item.ExplainerPotential,
			item.Timeliness,
			item.AudienceRelevance,
			time.Now().UTC(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
