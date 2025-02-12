package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chyiyaqing/gogpt/internal/config"
	"github.com/chyiyaqing/gogpt/internal/models"
)

type Message = models.Message
type ChatRequest = models.ChatRequest
type ChatResponse = models.ChatResponse

type Client struct {
	baseURL string
	client  *http.Client
	config  *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		client:  &http.Client{},
		config:  cfg,
	}
}

func (c *Client) buildOptions() map[string]any {
	return map[string]any{
		"temperature":    c.config.ModelParams.Temperature,
		"top_p":          c.config.ModelParams.TopP,
		"top_k":          c.config.ModelParams.TopK,
		"repeat_penalty": c.config.ModelParams.RepeatPenalty,
		"max_tokens":     c.config.ModelParams.MaxTokens,
	}
}

func (c *Client) ChatStream(ctx context.Context, model string, messages []Message, handler func(*ChatResponse)) error {
	reqBody := &ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
		Options:  c.buildOptions(),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request failed: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("read response failed: %v", err)
			}

			var chatResp ChatResponse
			if err := json.Unmarshal(line, &chatResp); err != nil {
				return fmt.Errorf("unmarshal response failed: %v", err)
			}

			if chatResp.Error != "" {
				return fmt.Errorf("api error: %s", chatResp.Error)
			}

			handler(&chatResp)

			if chatResp.Done {
				return nil
			}
		}
	}

}
