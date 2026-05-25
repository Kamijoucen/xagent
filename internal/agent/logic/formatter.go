package logic

import (
	"fmt"

	"github.com/admin/xagent/internal/agent/types"
	"github.com/admin/xagent/internal/appctx"
)

// FormatToolResult 将工具结果格式化成 LLM 可读文本。
func FormatToolResult(ctx *appctx.AppCtx, result types.ToolResult) string {
	_ = ctx
	if result.Success {
		return fmt.Sprintf("工具 %s 执行成功：%s", result.Name, result.Content)
	}
	return fmt.Sprintf("工具 %s 执行失败：%s", result.Name, result.Error)
}
