package models

type ArticleScore struct {
	ReelPotential      int `json:"reel_potential"`
	DeepDivePotential  int `json:"deep_dive_potential"`
	ExplainerPotential int `json:"explainer_potential"`
	Timeliness         int `json:"timeliness"`
	AudienceRelevance  int `json:"audience_relevance"`
	FinalScore         int `json:"final_score"`
}

// type ScoredNewsItem struct {
// 	NewsItem
// 	Score ArticleScore
// }

type ScoredNewsItem struct {
	NewsItem

	ReelPotential       float64 `json:"reel_potential"`
	DeepDivePotential   float64 `json:"deep_dive_potential"`
	DiscussionPotential float64 `json:"discussion_potential"`
	TutorialPotential   float64 `json:"tutorial_potential"`
	FinalScore          float64 `json:"final_score"`
	ExplainerPotential  float64 `json:"explainer_potential"`
	Timeliness          float64 `json:"timeliness"`
	AudienceRelevance   float64 `json:"audience_relevance"`
}
