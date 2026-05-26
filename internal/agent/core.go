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

	if ctx.TUI == nil {
		return fmt.Errorf("TUI 组件未初始化")
	}

	input := strings.TrimSpace(userInput)
	if input == "" {
		return fmt.Errorf("用户输入不能为空")
	}

	ctx.TUI.SetAgentState(types.AgentStateThinking)
	if ctx.SessionStore != nil {
		msg := types.Message{Role: types.RoleUser, Content: input, CreatedAt: time.Now()}
		if err := ctx.SessionStore.AppendMessage(components.DefaultSessionID, msg); err != nil {
			return fmt.Errorf("保存用户消息失败: %w", err)
		}
	}

	// 模拟思考时间
	select {
	case <-ctx.Context().Done():
		return context.Cause(ctx.Context())
	case <-time.After(120 * time.Millisecond):

	}

	ctx.TUI.SetAgentState(types.AgentStateResponding)

	var history []types.Message
	if ctx.SessionStore != nil {
		history = ctx.SessionStore.GetHistory(components.DefaultSessionID)
	}
	resp := ctx.LLMClient.StreamChat(ctx.Context(), history)

	var fullResponse strings.Builder
	for chunk := range resp {
		if chunk.Err != nil {
			fmt.Printf("LLM 流式输出错误: %v\n", chunk.Err)
			break
		}
		fullResponse.WriteString(chunk.Content)
		ctx.TUI.AppendAssistantMessage(chunk.Content)
	}
	ctx.TUI.SetAgentState(types.AgentStateIdle)
	if ctx.SessionStore != nil {
		msg := types.Message{Role: types.RoleAssistant, Content: fullResponse.String(), CreatedAt: time.Now()}
		if err := ctx.SessionStore.AppendMessage(components.DefaultSessionID, msg); err != nil {
			return fmt.Errorf("保存助手消息失败: %w", err)
		}
	}
	return nil
}
