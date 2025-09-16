package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamlucif3r/trikal/internal/config"
	"github.com/iamlucif3r/trikal/internal/ingest"
	"github.com/iamlucif3r/trikal/internal/logging"
	"github.com/iamlucif3r/trikal/internal/models"
	"github.com/iamlucif3r/trikal/internal/pipeline"
	"github.com/iamlucif3r/trikal/internal/ranking"
	"github.com/iamlucif3r/trikal/internal/scoring"
	"github.com/iamlucif3r/trikal/internal/store"
	"github.com/schollz/progressbar/v3"
)

func main() {
	cfg := config.MustLoad()
	logger := logging.New(cfg.Log)

	logger.Debug("starting application")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger.Debug("opening database", "url", cfg.DB.URL)
	db, err := store.Open(cfg.DB.URL)
	if err != nil {
		logger.Fatal("db open failed", "err", err)
	}
	defer func() {
		logger.Debug("closing database")
		db.Close()
	}()

	logger.Debug("creating HTTP client")
	httpc := ingest.NewHTTPClient(cfg)
	defer func() {
		logger.Debug("closing idle HTTP connections")
		httpc.CloseIdleConnections()
	}()

	logger.Debug("fetching items")
	items, err := ingest.FetchAll(ctx, httpc, cfg)
	if err != nil {
		logger.Error("ingest failed", "err", err)
	} else {
		logger.Info("items fetched", "count", len(items))
	}

	logger.Debug("normalizing items")
	items = pipeline.Normalize(items)
	logger.Info("items normalized", "count", len(items))

	logger.Debug("tagging items")
	items = pipeline.Tag(items)
	logger.Info("items tagged", "count", len(items))

	logger.Debug("deduplicating items in memory by URL")
	items = pipeline.DedupByURL(items)
	logger.Info("items deduplicated by URL", "count", len(items))

	logger.Debug("scoring items with LLM")
	bar := progressbar.NewOptions(len(items),
		progressbar.OptionSetDescription("Scoring articles..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: "-", BarStart: "[", BarEnd: "]"}),
	)
	var scoredItems []models.ScoredNewsItem
	for _, item := range items {
		scored, err := scoring.ScoreArticle(item)
		if err != nil {
			logger.Error("scoring failed", "title", item.Title, "err", err)
			continue
		}
		scoredItems = append(scoredItems, *scored)
		bar.Add(1)
	}

	logger.Info("items scored", "count", len(scoredItems))
	scoredItems = pipeline.DedupByURLScored(scoredItems)
	logger.Info("scored items deduplicated by URL", "count", len(scoredItems))

	logger.Debug("ranking items")
	ranked := ranking.Rank(scoredItems, ranking.ByComposite, 50, cfg)
	logger.Info("items ranked (composite)", "count", len(ranked))

	repo := store.NewNewsRepo(db)
	logger.Debug("upserting batch of items")
	if len(scoredItems) == 0 {
		logger.Warn("no items to upsert")
	} else {
		if err := repo.UpsertBatch(ctx, scoredItems); err != nil {
			logger.Error("batch upsert failed", "err", err, "count", len(scoredItems))
		} else {
			logger.Info("batch upsert succeeded", "count", len(scoredItems))
		}
	}

	logger.Info("ingest complete", "count", len(scoredItems))
	_ = time.Second
}
