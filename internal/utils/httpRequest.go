package utils

import (
	"fmt"
	"net/http"
)

func TriggerSarjan() error {
	url := "http://host.docker.internal:4444/generate"

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	fmt.Println("POST request successful!")
	return nil
}
