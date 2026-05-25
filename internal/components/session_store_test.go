package components

import (
	"testing"

	"github.com/admin/xagent/internal/agent/types"
)

func TestSessionStoreAppendAndGet(t *testing.T) {
	store, err := NewSessionStore("")
	if err != nil {
		t.Fatalf("创建会话存储失败: %v", err)
	}
	if err := store.AppendMessage(DefaultSessionID, types.Message{Role: types.RoleUser, Content: "hello"}); err != nil {
		t.Fatalf("追加消息失败: %v", err)
	}

	history := store.GetHistory(DefaultSessionID)
	if len(history) != 1 {
		t.Fatalf("历史数量错误: got %d", len(history))
	}
	if history[0].Content != "hello" {
		t.Fatalf("历史内容错误: got %q", history[0].Content)
	}
}
