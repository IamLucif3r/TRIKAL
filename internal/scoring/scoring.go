package scoring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/iamlucif3r/trikal/internal/models"
)

type ArticleScore struct {
	ReelPotential      float64 `json:"reel_potential"`
	DeepDivePotential  float64 `json:"deep_dive_potential"`
	ExplainerPotential float64 `json:"explainer_potential"`
	Timeliness         float64 `json:"timeliness"`
	AudienceRelevance  float64 `json:"audience_relevance"`
	FinalScore         float64 `json:"final_score"`
}

func ScoreArticle(news models.NewsItem) (*models.ScoredNewsItem, error) {
	prompt := fmt.Sprintf(`
You are scoring cybersecurity news for the hacker tabloid newsletter "pwnspectrum".

News:
Title: %s
Summary: %s

Rate the news on these criteria (0 = useless, 10 = perfect):

- reel_potential: Could this be a viral one-liner, meme, or shocking fact?
- deep_dive_potential: Does it have enough technical depth for a Medium article or YouTube walkthrough?
- explainer_potential: Can it be turned into a "How X Got Hacked" or educational explainer?
- timeliness: Is this breaking or currently trending?
- audience_relevance: Would this matter to bug bounty hunters, red/blue teams, CISOs, or security noobs?

Return ONLY valid JSON in this format:
{
  "reel_potential": 7,
  "deep_dive_potential": 9,
  "explainer_potential": 8,
  "timeliness": 10,
  "audience_relevance": 9
}
  Note: The score here are just representative. Do not copy them.
`, news.Title, news.Summary)

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	payload := map[string]any{
		"model":  os.Getenv("LLM_MODEL"),
		"prompt": prompt,
		"stream": false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama: %w", err)
	}
	defer resp.Body.Close()

	var respData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	completion, ok := respData["response"].(string)
	if !ok {
		return nil, fmt.Errorf("no completion text in Ollama response")
	}

	completion = strings.TrimSpace(completion)
	completion = strings.TrimSpace(completion)
	completion = strings.TrimPrefix(completion, "```json")
	completion = strings.TrimPrefix(completion, "```")
	completion = strings.TrimSuffix(completion, "```")
	completion = strings.TrimSpace(completion)

	var scores ArticleScore
	if err := json.Unmarshal([]byte(completion), &scores); err != nil {
		return nil, fmt.Errorf("failed to parse scores: %w\nraw: %s", err, completion)
	}

	scores.FinalScore = 0.3*scores.ReelPotential +
		0.25*scores.DeepDivePotential +
		0.25*scores.ExplainerPotential +
		0.1*scores.Timeliness +
		0.1*scores.AudienceRelevance

	return &models.ScoredNewsItem{
		NewsItem:            news,
		ReelPotential:       scores.ReelPotential,
		DeepDivePotential:   scores.DeepDivePotential,
		ExplainerPotential:  scores.ExplainerPotential,
		Timeliness:          scores.Timeliness,
		AudienceRelevance:   scores.AudienceRelevance,
		DiscussionPotential: scores.ExplainerPotential,
		TutorialPotential:   scores.DeepDivePotential,
		FinalScore:          scores.FinalScore,
	}, nil
}
