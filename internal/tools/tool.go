package tools

import (
	"context"

	"github.com/admin/xagent/internal/agent/types"
)

// Tool 是内置工具和插件工具的统一接口。
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args map[string]interface{}) (types.ToolResult, error)
}
