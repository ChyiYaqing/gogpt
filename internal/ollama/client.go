package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chyiyaqing/gogpt/internal/config"
	"github.com/chyiyaqing/gogpt/internal/models"
)

type Message = models.Message
type ChatRequest = models.ChatRequest
type ChatResponse = models.ChatResponse

type Client struct {
	baseURL    string
	httpClient *http.Client
	config     *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    60 * time.Second,
				DisableCompression: true,
			},
		},
		config: cfg,
	}
}

func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = baseURL
	return c
}

// buildOptions 构建模型参数
func (c *Client) buildOptions() map[string]any {
	return map[string]any{
		"temperature":    c.config.ModelParams.Temperature,
		"top_p":          c.config.ModelParams.TopP,
		"top_k":          c.config.ModelParams.TopK,
		"repeat_penalty": c.config.ModelParams.RepeatPenalty,
		"max_tokens":     c.config.ModelParams.MaxTokens,
	}
}

// Chat 非流式对话
func (c *Client) Chat(ctx context.Context, model string, messages []Message) (*models.ChatResponse, error) {
	reqBody := &ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options:  c.buildOptions(),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	var chatResp models.ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}
	if chatResp.Error != "" {
		return nil, fmt.Errorf("api error: %s", chatResp.Error)
	}

	return &chatResp, nil
}

// ChatStream implemets streaming chat functionality
func (c *Client) ChatStream(ctx context.Context, model string, messages []Message, handler ResponseHandler) error {
	options := StreamOptions{
		BufferSize:       c.config.StreamBufferSize,
		KeepFullResponse: c.config.KeepFullResponse,
		ErrorHandler: func(err error) {
			fmt.Printf("Stream error: %v\n", err)
		},
	}

	reqBody := ChatRequest{
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	streamReader := NewStreamReader(options, handler)
	return streamReader.Process(ctx, resp)

}

// ChatStreamWithOptions provides streaming with custom options
func (c *Client) ChatStreamWithOptions(ctx context.Context, model string, messages []Message, handler ResponseHandler, options StreamOptions) (string, error) {
	streamReader := NewStreamReader(options, handler)

	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
		Options:  c.buildOptions(),
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	if err := streamReader.Process(ctx, resp); err != nil {
		return "", err
	}

	if options.KeepFullResponse {
		return streamReader.GetFullResponse(), nil
	}
	return "", nil
}

// ListModels 获取可用模型列表
func (c *Client) ListModels(ctx context.Context) ([]models.ModelInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/models", nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	var models []models.ModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return models, nil
}

// GetModels 获取特定模型信息
func (c *Client) GetModel(ctx context.Context, name string) (*models.ModelInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/show/%s", c.baseURL, name), nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	var model models.ModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	return &model, nil
}

// Health checks if the Ollama server is running
func (c *Client) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/version", c.baseURL), nil)
	if err != nil {
		return fmt.Errorf("create health request failed: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server unhealthy, status: %d", resp.StatusCode)
	}

	return nil
}

// SaveConversation 保存对话
func (c *Client) SaveConversation(conv *models.Conversation, filepath string) error {
	return conv.SaveToFile(filepath)
}

// LoadConversation 加载对话
func (c *Client) LoadConversation(filepath string) (*models.Conversation, error) {
	return models.LoadConversationFromFile(filepath)
}
