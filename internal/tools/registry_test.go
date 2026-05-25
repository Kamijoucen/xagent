package tools

import (
	"context"
	"testing"

	"github.com/admin/xagent/internal/agent/types"
)

type fakeTool struct{}

func (fakeTool) Name() string        { return "fake" }
func (fakeTool) Description() string { return "fake tool" }
func (fakeTool) Execute(ctx context.Context, args map[string]interface{}) (types.ToolResult, error) {
	_ = ctx
	_ = args
	return types.ToolResult{Name: "fake", Content: "ok", Success: true}, nil
}

func TestRegistryRegisterGetList(t *testing.T) {
	registry := NewRegistry()
	if err := registry.Register("fake", fakeTool{}); err != nil {
		t.Fatalf("注册工具失败: %v", err)
	}
	if _, ok := registry.Get("fake"); !ok {
		t.Fatal("应能获取已注册工具")
	}
	if got := len(registry.List()); got != 1 {
		t.Fatalf("工具数量错误: got %d", got)
	}
}
