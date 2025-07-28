package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	OllamaAPIURL = "http://localhost:11434/api/generate"
)

func QueryOllamaAPI(prompt string) (float64, error) {

	payload := map[string]any{
		"model":  os.Getenv("LLM_MODEL"),
		"prompt": prompt,
		"stream": false,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(OllamaAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("non-200 response from Ollama API: %d %s", resp.StatusCode, string(bodyBytes))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return 0, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var completion string
	if c, ok := responseData["completion"].(string); ok {
		completion = c
	} else if c, ok := responseData["response"].(string); ok {
		completion = c
	} else {
		return 0, fmt.Errorf("could not find completion text in response")
	}

	cleaned := strings.TrimSpace(completion)
	fields := strings.Fields(cleaned)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty completion received")
	}
	score, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse score: %v", err)
	}
	return score, nil
}
