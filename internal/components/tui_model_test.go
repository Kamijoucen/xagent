package components

import (
	"regexp"
	"strings"
	"testing"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func TestNewTUIModelShowsFullPlaceholder(t *testing.T) {
	model := NewTUIModel()

	view := ansiEscapeRe.ReplaceAllString(model.input.View(), "")
	if !strings.Contains(view, "输入消息，按 Enter 发送") {
		t.Fatalf("placeholder 被截断: %q", view)
	}
}

func TestNewTUIModelShowsStartupBanner(t *testing.T) {
	model := NewTUIModel()

	view := model.View().Content
	if !strings.Contains(view, "__  __") || !strings.Contains(view, "|___/") {
		t.Fatalf("启动字符画未显示: %q", view)
	}
}
