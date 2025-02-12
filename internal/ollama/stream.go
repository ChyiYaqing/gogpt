package ollama

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// StreamProcessor 处理流式响应的接口
type StreamProcessor interface {
	Process(context.Context, *http.Response) error
}

// ResponseHandler 定义响应处理函数类型
type ResponseHandler func(*ChatResponse)

// StreamOptions 流处理选项
type StreamOptions struct {
	// 是否保留完整响应
	KeepFullResponse bool
	// 自定义分隔符
	Delimiter string
	// 缓冲区大小
	BufferSize int
	// 错误处理函数
	ErrorHandler func(error)
}

// DefaultStreamOptions 默认流处理选项
var DefaultStreamOptions = StreamOptions{
	KeepFullResponse: true,
	Delimiter:        "\n",
	BufferSize:       4096,
	ErrorHandler:     func(err error) { fmt.Printf("Stream error: %v\n", err) },
}

// StreamReader 流式读取器
type StreamReader struct {
	reader    *bufio.Reader
	options   StreamOptions
	handler   ResponseHandler
	buffer    strings.Builder
	mutex     sync.Mutex
	lastError error
}

// NewStreamReader 创建新的流式读取器
func NewStreamReader(options StreamOptions, handler ResponseHandler) *StreamReader {
	if options.BufferSize == 0 {
		options.BufferSize = 4096
	}
	return &StreamReader{
		options: options,
		handler: handler,
	}
}

// Process 处理流式响应
func (sr *StreamReader) Process(ctx context.Context, resp *http.Response) error {
	sr.reader = bufio.NewReaderSize(resp.Body, sr.options.BufferSize)

	// 创建错误通道
	errChan := make(chan error, 1)

	// 启动处理协程
	go func() {
		errChan <- sr.processStream(ctx)
	}()

	// 等待处理完成或上下文取消
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (sr *StreamReader) processStream(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := sr.reader.ReadBytes('\n')
			if err == io.EOF {
				return sr.handleFinalResponse()
			}
			if err != nil {
				return fmt.Errorf("read stream failed: %v", err)
			}

			if err := sr.handleResponse(line); err != nil {
				if sr.options.ErrorHandler != nil {
					sr.options.ErrorHandler(err)
				}
				continue
			}
		}
	}
}

// handleResponse 处理单条响应
func (sr *StreamReader) handleResponse(line []byte) error {
	// 跳过空行
	if len(strings.TrimSpace(string(line))) == 0 {
		return nil
	}

	var response ChatResponse
	if err := json.Unmarshal(line, &response); err != nil {
		return fmt.Errorf("unmarshal response failed: %v", err)
	}

	// 检查API错误
	if response.Error != "" {
		return fmt.Errorf("api error: %s", response.Error)
	}

	// 如果需要保留完整响应，追加到buffer
	if sr.options.KeepFullResponse {
		sr.appendToBuffer(response.Message.Content)
	}

	// 调用处理函数
	if sr.handler != nil {
		sr.handler(&response)
	}

	return nil
}

// appendToBuffer 追加内容到缓冲区
func (sr *StreamReader) appendToBuffer(content string) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	sr.buffer.WriteString(content)
}

// GetFullResponse 获取完整响应
func (sr *StreamReader) GetFullResponse() string {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	return sr.buffer.String()
}

// handleFinalResponse 处理最后一条响应
func (sr *StreamReader) handleFinalResponse() error {
	if sr.lastError != nil {
		return sr.lastError
	}
	return nil
}

// ResetBuffer 重置缓冲区
func (sr *StreamReader) ResetBuffer() {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	sr.buffer.Reset()
}
