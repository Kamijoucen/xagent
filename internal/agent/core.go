package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/admin/xagent/internal/agent/types"
	"github.com/admin/xagent/internal/appctx"
	"github.com/admin/xagent/internal/components"
)

// RunReActLoop 是 MVP 版编排层：只做本地假响应，保留后续替换为 ReAct 的位置。
func RunReActLoop(ctx *appctx.AppCtx, userInput string) error {
	if ctx == nil {
		return fmt.Errorf("AppCtx 不能为空")
	}
	input := strings.TrimSpace(userInput)
	if input == "" {
		return fmt.Errorf("用户输入不能为空")
	}

	if ctx.TUI != nil {
		ctx.TUI.SetAgentState(types.AgentStateThinking)
	}
	if ctx.SessionStore != nil {
		if err := ctx.SessionStore.AppendMessage(components.DefaultSessionID, types.Message{Role: types.RoleUser, Content: input, CreatedAt: time.Now()}); err != nil {
			return fmt.Errorf("保存用户消息失败: %w", err)
		}
	}

	select {
	case <-ctx.Context().Done():
		return context.Cause(ctx.Context())
	case <-time.After(120 * time.Millisecond):
	}

	if ctx.TUI != nil {
		ctx.TUI.SetAgentState(types.AgentStateResponding)
	}
	response := "收到：" + input + "\n\n这是第一版 MVP 的本地响应。后续可以在这里接入 LLM 流式输出、工具调用和会话持久化。"
	if ctx.TUI != nil {
		ctx.TUI.AppendAssistantMessage(response)
		ctx.TUI.SetAgentState(types.AgentStateIdle)
	}
	if ctx.SessionStore != nil {
		if err := ctx.SessionStore.AppendMessage(components.DefaultSessionID, types.Message{Role: types.RoleAssistant, Content: response, CreatedAt: time.Now()}); err != nil {
			return fmt.Errorf("保存助手消息失败: %w", err)
		}
	}
	return nil
}
