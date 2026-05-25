package logic

import (
	"fmt"

	"golang.org/x/sync/semaphore"

	"github.com/admin/xagent/internal/agent/types"
	"github.com/admin/xagent/internal/appctx"
)

// ExecuteToolChain 执行工具链。MVP 保留 semaphore 并发边界，但只做最小串行执行。
func ExecuteToolChain(ctx *appctx.AppCtx, calls []types.ToolCall) ([]types.ToolResult, error) {
	if len(calls) == 0 {
		return nil, nil
	}
	if ctx == nil || ctx.ToolRegistry == nil {
		return nil, fmt.Errorf("工具注册表未初始化")
	}

	limit := semaphore.NewWeighted(5)
	results := make([]types.ToolResult, 0, len(calls))
	for _, call := range calls {
		if err := limit.Acquire(ctx.Context(), 1); err != nil {
			return results, err
		}
		tool, ok := ctx.ToolRegistry.Get(call.Name)
		if !ok {
			results = append(results, types.ToolResult{Name: call.Name, Success: false, Error: "工具未注册"})
			limit.Release(1)
			continue
		}
		result, err := tool.Execute(ctx.Context(), call.Args)
		if err != nil {
			result = types.ToolResult{Name: call.Name, Success: false, Error: err.Error()}
		}
		results = append(results, result)
		limit.Release(1)
	}
	return results, nil
}
