package components

import (
	"strings"
	"testing"
)

func TestNewTUIModelShowsFullPlaceholder(t *testing.T) {
	model := NewTUIModel()

	view := model.input.View()
	if !strings.Contains(view, "输入消息，按 Enter 发送") {
		t.Fatalf("placeholder 被截断: %q", view)
	}
}

func TestNewTUIModelShowsStartupBanner(t *testing.T) {
	model := NewTUIModel()

	view := model.View()
	if !strings.Contains(view, "__  __") || !strings.Contains(view, "|___/") {
		t.Fatalf("启动字符画未显示: %q", view)
	}
}
