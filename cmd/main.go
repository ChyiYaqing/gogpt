package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Message 定义消息结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 定义请求结构体
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse 定义响应结构体
type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
	Error   string  `json:"error,omitempty"`
}

// OllamaClient Ollama 客户端结构体
type OllamaClient struct {
	baseURL string
	client  *http.Client
}

// NewOllamaClient 创建新的Ollama客户端
func NewOllamaClient(baseURL string) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Chat 发送对话请求
func (c *OllamaClient) Chat(modelName string, messages []Message) (*ChatResponse, error) {
	// 构件请求体
	reqBody := &ChatRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", c.baseURL+"/api/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %v", err)
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %v", err)
	}

	if chatResp.Error != "" {
		return nil, fmt.Errorf("api error: %s", chatResp.Error)
	}

	return &chatResp, nil
}

func main() {
	// 创建客户端
	client := NewOllamaClient("http://localhost:11434")

	// 准备对话消息
	messages := []Message{
		{
			Role:    "user",
			Content: "你好，请介绍一下自己",
		},
	}

	// 发送请求
	resp, err := client.Chat("deepseek-r1:8b", messages)
	if err != nil {
		fmt.Printf("Chat failed: %v\n", err)
		return
	}

	// 打印响应
	fmt.Printf("Assistant: %s\n", resp.Message.Content)
}
