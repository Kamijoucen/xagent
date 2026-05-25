package components

import (
	"context"
	"io"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/tmaxmax/go-sse"

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

// LLMClient 是 go-openai 的业务组件封装。第一版不发起真实请求。
type LLMClient struct {
	client *openai.Client
	cfg    LLMClientConfig
}

// NewLLMClient 创建 LLM 客户端占位实现。
func NewLLMClient(cfg LLMClientConfig) *LLMClient {
	openAIConfig := openai.DefaultConfig(cfg.APIKey)
	if cfg.APIBase != "" {
		openAIConfig.BaseURL = cfg.APIBase
	}
	return &LLMClient{client: openai.NewClientWithConfig(openAIConfig), cfg: cfg}
}

// StreamChat 第一版固定走本地假流式输出，后续可替换为真实 LLM。
func (c *LLMClient) StreamChat(ctx context.Context, prompt string) <-chan types.StreamChunk {
	return c.MockStreamChat(ctx, prompt)
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
			out <- types.StreamChunk{Content: "收到：" + strings.TrimSpace(prompt), Done: true}
		}
	}()
	return out
}

// DecodeSSEPayloads 是预留的 SSE 解析 helper，当前只用于保持依赖和边界清晰。
func DecodeSSEPayloads(reader io.Reader) ([]string, error) {
	var payloads []string
	var readErr error
	events := sse.Read(reader, nil)
	events(func(event sse.Event, err error) bool {
		if err != nil {
			readErr = err
			return false
		}
		payloads = append(payloads, event.Data)
		return true
	})
	return payloads, readErr
}

// Shutdown 释放 LLM 客户端资源。当前实现无持久连接。
func (c *LLMClient) Shutdown(ctx context.Context) error {
	_ = ctx
	return nil
}
