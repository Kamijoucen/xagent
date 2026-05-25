package logic

import (
	"strings"

	"github.com/admin/xagent/internal/agent/types"
	"github.com/admin/xagent/internal/appctx"
)

// BuildSystemPrompt 根据工具列表构建系统提示词。MVP 只返回简洁文本。
func BuildSystemPrompt(ctx *appctx.AppCtx, tools []types.ToolSpec) string {
	_ = ctx
	if len(tools) == 0 {
		return "你是 agent-cli 的 MVP 本地助手。当前没有启用工具。"
	}

	var builder strings.Builder
	builder.WriteString("你是 agent-cli 的 MVP 本地助手。可用工具：")
	for _, tool := range tools {
		builder.WriteString("\n- ")
		builder.WriteString(tool.Name)
		if tool.Description != "" {
			builder.WriteString(": ")
			builder.WriteString(tool.Description)
		}
	}
	return builder.String()
}

// BuildUserPrompt 组装用户输入和历史上下文。
func BuildUserPrompt(ctx *appctx.AppCtx, input string, history []types.Message) string {
	_ = ctx
	var builder strings.Builder
	for _, message := range history {
		builder.WriteString(string(message.Role))
		builder.WriteString(": ")
		builder.WriteString(message.Content)
		builder.WriteString("\n")
	}
	builder.WriteString("user: ")
	builder.WriteString(strings.TrimSpace(input))
	return builder.String()
}
