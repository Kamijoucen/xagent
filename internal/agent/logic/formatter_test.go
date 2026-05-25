package logic

import (
	"strings"
	"testing"

	"github.com/admin/xagent/internal/agent/types"
)

func TestFormatToolResult(t *testing.T) {
	text := FormatToolResult(nil, types.ToolResult{Name: "demo", Content: "ok", Success: true})
	if !strings.Contains(text, "demo") || !strings.Contains(text, "ok") {
		t.Fatalf("格式化结果不符合预期: %q", text)
	}
}
