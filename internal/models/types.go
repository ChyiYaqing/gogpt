package models

import (
	"encoding/json"
	"os"
	"time"
)

// Message 定义消息结构体
type Message struct {
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time,omitempty"`
}

// ChatRequest 定义请求结构体
type ChatRequest struct {
	Model    string         `json:"model"`
	Messages []Message      `json:"messages"`
	Stream   bool           `json:"stream"`
	Options  map[string]any `json:"options,omitempty"`
	Template string         `json:"template,omitempty"` // 可选的提示模版
	Format   string         `json:"format,omitempty`    // 响应格式
	Context  []int          `json:"context,omitempty"`  // 上下文窗口
}

// ChatResponse 定义响应结构体
type ChatResponse struct {
	Model         string      `json:"model"`
	Message       Message     `json:"message"`
	Done          bool        `json:"done"`
	TotalDuration float64     `json:"total_duration,omitempty"`
	LoadDuration  float64     `json:"load_duration,omitempty"`
	PromptEval    EvalMetrics `json:"prompt_eval_duration, omitempty"`
	EvalCount     int         `json:"eval_count,omitempty"`
	Error         string      `json:"error,omitempty"`
	Context       []int       `json:"context,omitempty"`
}

// EvalMetrics 定义评估指标
type EvalMetrics struct {
	Duration float64 `json:"duration"`
	Tokens   int     `json:"tokens"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	Name     string       `json:"name"`
	Modified time.Time    `json:"modified"`
	Size     int64        `json:"size"`
	Digest   string       `json:"digest"`
	Details  ModelDetails `json:"details,omitempty"`
}

// ModelDetails 模型详细信息
type ModelDetails struct {
	Format           string   `json:"format"`
	Family           string   `json:"family"`
	Families         []string `json:"families"`
	ParameterSize    string   `json:"parameter_size"`
	QuantizationType string   `json:"quantization_type"`
}

// ConversationMeta 对话元数据
type ConversationMeta struct {
	ID         string    `json:"id"`
	Model      string    `json:"model"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	TurnCount  int       `json:"turn_count"`
	TokenCount int       `json:"token_count"`
}

// Conversation 完整对话记录
type Conversation struct {
	Meta     ConversationMeta `json:"meta"`
	Messages []Message        `json:"message"`
}

// SaveToFile 保存对话到文件
func (c *Conversation) SaveToFile(filepath string) error {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

// LoadFormFile 从文件加载对话
func LoadConversationFromFile(filepath string) (*Conversation, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var conv Conversation
	if err := json.Unmarshal(data, &conv); err != nil {
		return nil, err
	}
	return &conv, nil
}
