package types

import "time"

// Role 表示对话消息角色。
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message 是 MVP 中跨组件共享的对话消息。
type Message struct {
	Role      Role
	Content   string
	CreatedAt time.Time
}

// AgentState 表示 Agent 当前的用户可见状态。
type AgentState string

const (
	AgentStateIdle        AgentState = "idle"
	AgentStateThinking    AgentState = "thinking"
	AgentStateResponding  AgentState = "responding"
	AgentStateToolCalling AgentState = "tool_calling"
	AgentStateError       AgentState = "error"
)

// StreamChunk 是 LLM 流式输出的轻量占位类型。
type StreamChunk struct {
	Content string
	Done    bool
	Err     error
}

// ToolSpec 描述一个可用工具，供 prompt 构建使用。
type ToolSpec struct {
	Name        string
	Description string
}

// ToolCall 表示一次工具调用请求。
type ToolCall struct {
	Name string
	Args map[string]interface{}
}

// ToolResult 表示工具执行结果。
type ToolResult struct {
	Name    string
	Content string
	Success bool
	Error   string
}

// LLMOutput 是解析后的 LLM 输出。
type LLMOutput struct {
	Text         string
	ToolCalls    []ToolCall
	HasToolCalls bool
}

// PluginResult 是插件系统预留的返回类型。
type PluginResult struct {
	Content string
	Success bool
	Error   string
}
