package types

type Config struct {
	DatabaseURL string `json:"db_url"`
	OllamaURL   string `json:"ollama_url"`
}
