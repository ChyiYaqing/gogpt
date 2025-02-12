package models

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
	Options map[string]any `json:"options,omitempty"`
}

// ChatResponse 定义响应结构体
type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
	Error   string  `json:"error,omitempty"`
}
