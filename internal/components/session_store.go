package components

import (
	"context"
	"sync"

	bolt "go.etcd.io/bbolt"

	"github.com/admin/xagent/internal/agent/types"
)

// DefaultSessionID 是 MVP 单会话 ID。
const DefaultSessionID = "default"

// SessionStore 第一版使用内存消息列表，保留 bbolt 字段作为后续扩展点。
type SessionStore struct {
	mu       sync.RWMutex
	dbPath   string
	db       *bolt.DB
	messages map[string][]types.Message
}

// NewSessionStore 创建内存会话存储。
func NewSessionStore(dbPath string) (*SessionStore, error) {
	return &SessionStore{
		dbPath: dbPath,
		messages: map[string][]types.Message{
			DefaultSessionID: {},
		},
	}, nil
}

// GetHistory 返回指定会话历史的副本。
func (s *SessionStore) GetHistory(sessionID string) []types.Message {
	if sessionID == "" {
		sessionID = DefaultSessionID
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	history := s.messages[sessionID]
	clone := make([]types.Message, len(history))
	copy(clone, history)
	return clone
}

// AppendMessage 追加会话消息。
func (s *SessionStore) AppendMessage(sessionID string, msg types.Message) error {
	if sessionID == "" {
		sessionID = DefaultSessionID
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[sessionID] = append(s.messages[sessionID], msg)
	return nil
}

// Shutdown 关闭会话存储。当前只有可选 bbolt 句柄。
func (s *SessionStore) Shutdown(ctx context.Context) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}
