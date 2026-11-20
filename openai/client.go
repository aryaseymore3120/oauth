// Package openai provides a small client for the OpenAI Responses API.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	key string
	http *http.Client
	endpoint string
}

// NewClient reads the API key from the environment.
func NewClient() (*Client, error) {
	key := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if key == "" {
		return nil, errors.New("OPENAI_API_KEY is required")
	}
	return &Client{key: key, http: &http.Client{
		Timeout: 60 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, endpoint: "https://api.openai.com/v1/responses"}, nil
}

// Generate sends a text prompt and returns the response text.
func (c *Client) Generate(ctx context.Context, model, input string) (string, error) {
	if strings.TrimSpace(model) == "" || strings.TrimSpace(input) == "" {
		return "", errors.New("model and input are required")
	}
	body, err := json.Marshal(map[string]interface{}{"model": model, "input": input, "store": false})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("OpenAI API returned HTTP %d", res.StatusCode)
	}
	var result struct {
		Status string `json:"status"`
		Output []struct {
			Type string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Status != "completed" {
		return "", fmt.Errorf("OpenAI response not completed: %s", result.Status)
	}
	var text strings.Builder
	for _, item := range result.Output {
		if item.Type != "message" {
			continue
		}
		for _, part := range item.Content {
			if part.Type == "output_text" {
				text.WriteString(part.Text)
			}
		}
	}
	if text.Len() == 0 {
		return "", errors.New("OpenAI returned no output text")
	}
	return text.String(), nil
}
