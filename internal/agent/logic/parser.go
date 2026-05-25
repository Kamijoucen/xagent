package logic

import (
	"strings"

	"github.com/admin/xagent/internal/agent/types"
	"github.com/admin/xagent/internal/appctx"
)

// ParseLLMResponse 解析 LLM 输出。MVP 仅支持纯文本。
func ParseLLMResponse(ctx *appctx.AppCtx, raw string) (*types.LLMOutput, error) {
	_ = ctx
	return &types.LLMOutput{Text: strings.TrimSpace(raw)}, nil
}

// ExtractToolCalls 从文本中提取工具调用。MVP 暂不启用工具调用。
func ExtractToolCalls(ctx *appctx.AppCtx, content string) ([]types.ToolCall, error) {
	_ = ctx
	_ = content
	return nil, nil
}
