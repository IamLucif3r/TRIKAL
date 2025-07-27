package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/iamlucif3r/trikal/internal/database"
	"github.com/iamlucif3r/trikal/internal/types"
)

func SendAlertToDiscord(webhookURL string) error {
	db := database.DB
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	var maxScore float64
	err := db.QueryRow(`SELECT MAX(llm_score) FROM articles`).Scan(&maxScore)
	if err != nil {
		return fmt.Errorf("failed to query max llm_score: %v", err)
	}

	rows, err := db.Query(
		`SELECT title, description, link, source, published_at FROM articles WHERE llm_score = $1`, maxScore)
	if err != nil {
		return fmt.Errorf("failed to query articles: %v", err)
	}
	defer rows.Close()

	var embeds []types.DiscordEmbed
	for rows.Next() {
		var title, description, link, source, publishedAt string
		if err := rows.Scan(&title, &description, &link, &source, &publishedAt); err != nil {
			log.Printf("failed to scan article: %v", err)
			continue
		}

		desc := description
		if len(desc) > 500 {
			desc = desc[:497] + "..."
		}
		if len(title) > 256 {
			title = title[:253] + "..."
		}

		embed := types.DiscordEmbed{
			Title:       title,
			Description: desc,
			URL:         link,
			Color:       0x4287f5,
			Fields: []types.DiscordEmbedField{
				{Name: "Source", Value: source, Inline: true},
				{Name: "Published At", Value: publishedAt, Inline: true},
				{Name: "LLM Score", Value: fmt.Sprintf("%.2f", maxScore), Inline: true},
			},
		}

		embeds = append(embeds, embed)
	}

	if len(embeds) == 0 {
		log.Println("[INFO] No top-score articles to send.")
		return nil
	}

	batchSize := 10
	for i := 0; i < len(embeds); i += batchSize {
		end := i + batchSize
		if end > len(embeds) {
			end = len(embeds)
		}
		batch := embeds[i:end]

		payload := map[string]interface{}{
			"embeds":  batch,
			"content": fmt.Sprintf("🔥 **Top Cybersecurity News (LLM Score: %.2f)** 🔥", maxScore),
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Could not marshal payload: %v", err)
			continue
		}

		req, err := http.NewRequest("POST", webhookURL, strings.NewReader(string(payloadBytes)))
		if err != nil {
			log.Printf("Could not create POST request: %v", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Could not post to Discord: %v", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Printf("Discord webhook returned status %d", resp.StatusCode)
			continue
		}
		log.Printf("[SUCCESS] Sent a batch of %d top-score article(s) to Discord.", len(batch))
	}

	return nil
}
