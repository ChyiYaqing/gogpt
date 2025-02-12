package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chyiyaqing/gogpt/internal/config"
	"github.com/chyiyaqing/gogpt/internal/ollama"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 创建客户端
	client := ollama.NewClient(cfg)

	// 创建上下文和取消函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Received interrupt signal, exiting...")
		cancel()
	}()

	// 创建对话历史
	conversation := []ollama.Message{}

	// 启动交互对话
	fmt.Println("Welcome to GPT chatbot! Type 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if strings.ToLower(input) == "exit" {
			break
		}

		// 添加用户消息对话历史
		conversation = append(conversation, ollama.Message{
			Role:    "user",
			Content: input,
		})

		// 创建带超时的上下文
		ctxWithTimeout, cancelTimeout := context.WithTimeout(ctx, cfg.Timeout)

		// 处理流式响应
		fmt.Print("\nBot: ")
		response := ""
		err := client.ChatStream(ctxWithTimeout, cfg.Model, conversation, func(resp *ollama.ChatResponse) {
			if resp.Message.Content != "" {
				fmt.Print(resp.Message.Content)
				response += resp.Message.Content
			}
		})

		cancelTimeout()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// 添加助手回复到对话历史
		conversation = append(conversation, ollama.Message{
			Role:    "assistant",
			Content: response,
		})

		// 如果对话历史过长，保留最近的n轮对话
		if len(conversation) > cfg.MaxConversationTurns*2 {
			conversation = conversation[2:]
		}
	}
}
