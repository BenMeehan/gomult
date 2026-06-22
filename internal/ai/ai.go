package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const deepseekAPI = "https://api.deepseek.com/v1/chat/completions"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

func (c *Client) Analyze(code, language, input, output, status string) (string, error) {
	var prompt string
	if status == "ok" {
		prompt = fmt.Sprintf(
			"Analyze the following %s code briefly. Comment on code quality, potential issues, and suggestions for improvement.\n\nCode:\n```%s\n%s\n```\n\nInput: %s\n\nOutput:\n%s",
			language, language, code, input, output,
		)
	} else {
		prompt = fmt.Sprintf(
			"The following %s code produced an error. Identify the bug and suggest a fix.\n\nCode:\n```%s\n%s\n```\n\nInput: %s\n\n%s",
			language, language, code, input, output,
		)
	}

	reqBody := chatRequest{
		Model: "deepseek-chat",
		Messages: []message{
			{Role: "system", Content: "You are a code analysis assistant. Be concise and direct. Use bullet points when helpful."},
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, deepseekAPI, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("api error %d: %s", resp.StatusCode, string(errBody))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return chatResp.Choices[0].Message.Content, nil
}
