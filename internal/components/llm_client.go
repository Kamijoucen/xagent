package components

import (
	"context"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/admin/xagent/internal/agent/types"
)

// LLMClientConfig 保存 LLM 客户端配置。
type LLMClientConfig struct {
	APIKey      string
	APIBase     string
	Model       string
	MaxTokens   int
	Temperature float32
	MockMode    bool
}

// LLMClient 是 openai-go 官方库的业务组件封装。
type LLMClient struct {
	client openai.Client
	cfg    LLMClientConfig
}

// NewLLMClient 创建 LLM 客户端。
func NewLLMClient(cfg LLMClientConfig) *LLMClient {
	opts := []option.RequestOption{option.WithAPIKey(cfg.APIKey)}
	if cfg.APIBase != "" {
		opts = append(opts, option.WithBaseURL(cfg.APIBase))
	}
	return &LLMClient{client: openai.NewClient(opts...), cfg: cfg}
}

// StreamChat 根据 MockMode 选择走本地假输出或真实 LLM 流式请求。
// history 为多轮对话上下文，mock 模式下只取最后一条 user message。
func (c *LLMClient) StreamChat(ctx context.Context, history []types.Message) <-chan types.StreamChunk {
	if c.cfg.MockMode {
		prompt := ""
		for i := len(history) - 1; i >= 0; i-- {
			if history[i].Role == types.RoleUser {
				prompt = history[i].Content
				break
			}
		}
		return c.MockStreamChat(ctx, prompt)
	}
	return c.RealStreamChat(ctx, history)
}

// RealStreamChat 使用官方 openai-go 库发起真实流式 Chat Completion 请求。
func (c *LLMClient) RealStreamChat(ctx context.Context, history []types.Message) <-chan types.StreamChunk {
	out := make(chan types.StreamChunk)
	go func() {
		defer close(out)

		messages := make([]openai.ChatCompletionMessageParamUnion, 0, len(history))
		for _, msg := range history {
			switch msg.Role {
			case types.RoleSystem:
				messages = append(messages, openai.SystemMessage(msg.Content))
			case types.RoleUser:
				messages = append(messages, openai.UserMessage(msg.Content))
			case types.RoleAssistant:
				messages = append(messages, openai.AssistantMessage(msg.Content))
			}
		}

		params := openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    c.cfg.Model,
		}
		if c.cfg.MaxTokens > 0 {
			params.MaxTokens = openai.Int(int64(c.cfg.MaxTokens))
		}
		if c.cfg.Temperature >= 0 {
			params.Temperature = openai.Float(float64(c.cfg.Temperature))
		}

		stream := c.client.Chat.Completions.NewStreaming(ctx, params)
		for stream.Next() {
			event := stream.Current()
			if len(event.Choices) == 0 {
				continue
			}
			delta := event.Choices[0].Delta.Content
			if delta == "" {
				continue
			}
			select {
			case out <- types.StreamChunk{Content: delta}:
			case <-ctx.Done():
				out <- types.StreamChunk{Err: ctx.Err(), Done: true}
				return
			}
		}

		if err := stream.Err(); err != nil {
			select {
			case out <- types.StreamChunk{Err: err, Done: true}:
			case <-ctx.Done():
			}
			return
		}

		select {
		case out <- types.StreamChunk{Done: true}:
		case <-ctx.Done():
		}
	}()
	return out
}

// MockStreamChat 返回一条简单的本地响应。
func (c *LLMClient) MockStreamChat(ctx context.Context, prompt string) <-chan types.StreamChunk {
	out := make(chan types.StreamChunk, 1)
	go func() {
		defer close(out)
		select {
		case <-ctx.Done():
			out <- types.StreamChunk{Err: ctx.Err(), Done: true}
		case <-time.After(60 * time.Millisecond):
			out <- types.StreamChunk{Content: "收到了：" + strings.TrimSpace(prompt), Done: true}
		}
	}()
	return out
}

// Shutdown 释放 LLM 客户端资源。当前实现无持久连接。
func (c *LLMClient) Shutdown(ctx context.Context) error {
	_ = ctx
	return nil
}
