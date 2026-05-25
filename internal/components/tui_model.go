package components

import (
	"errors"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/admin/xagent/internal/agent/types"
)

// SubmitHandler 是 TUI 提交用户输入时调用的启动器注入函数。
type SubmitHandler func(input string) tea.Cmd

// TUIModel 是 Bubble Tea Model，也是 ApplicationContext 持有的 UI Component。
type TUIModel struct {
	messages []types.Message
	input    textinput.Model
	viewport viewport.Model
	state    types.AgentState
	errText  string
	width    int
	height   int

	mu            sync.RWMutex
	sender        func(tea.Msg)
	submitHandler SubmitHandler
}

// AgentStateMsg 更新 Agent 可见状态。
type AgentStateMsg struct {
	State types.AgentState
}

// AssistantMessageMsg 追加一条助手消息。
type AssistantMessageMsg struct {
	Content string
}

// ErrorMsg 表示异步任务错误。
type ErrorMsg struct {
	Err string
}

// NewTUIModel 创建带欢迎消息的 TUI 模型。
func NewTUIModel() *TUIModel {
	input := textinput.New()
	input.Placeholder = "输入消息，按 Enter 发送"
	input.Prompt = "> "
	input.Focus()
	input.CharLimit = 4096
	input.Width = 80

	vp := viewport.New(80, 18)
	model := &TUIModel{
		messages: []types.Message{
			{Role: types.RoleSystem, Content: "欢迎使用 agent-cli MVP。第一版保留清晰架构，只提供本地假交互。", CreatedAt: time.Now()},
			{Role: types.RoleAssistant, Content: "你可以直接输入任意文本，我会返回一条本地响应。", CreatedAt: time.Now()},
		},
		input:    input,
		viewport: vp,
		state:    types.AgentStateIdle,
		width:    80,
		height:   24,
	}
	model.syncViewport()
	return model
}

// Init 初始化 Bubble Tea Model。
func (m *TUIModel) Init() tea.Cmd {
	return textinput.Blink
}

// SetSender 注入 Bubble Tea Program.Send，用于后台任务安全更新 UI。
func (m *TUIModel) SetSender(sender func(tea.Msg)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sender = sender
}

// SetSubmitHandler 注入用户输入提交后的启动逻辑。
func (m *TUIModel) SetSubmitHandler(handler SubmitHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.submitHandler = handler
}

// SetAgentState 通过 Bubble Tea 消息更新状态。
func (m *TUIModel) SetAgentState(state types.AgentState) {
	m.emit(AgentStateMsg{State: state})
}

// AppendAssistantMessage 追加助手消息。
func (m *TUIModel) AppendAssistantMessage(content string) {
	m.emit(AssistantMessageMsg{Content: content})
}

// ReportError 将错误展示到状态栏。
func (m *TUIModel) ReportError(err error) {
	if err == nil {
		return
	}
	m.emit(ErrorMsg{Err: err.Error()})
}

func (m *TUIModel) emit(msg tea.Msg) {
	m.mu.RLock()
	sender := m.sender
	m.mu.RUnlock()
	if sender != nil {
		sender(msg)
		return
	}
	m.applyMessage(msg)
}

func (m *TUIModel) currentSubmitHandler() SubmitHandler {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.submitHandler
}

func (m *TUIModel) applyMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case AgentStateMsg:
		m.state = msg.State
		if msg.State != types.AgentStateError {
			m.errText = ""
		}
	case AssistantMessageMsg:
		m.messages = append(m.messages, types.Message{Role: types.RoleAssistant, Content: msg.Content, CreatedAt: time.Now()})
		m.syncViewport()
	case ErrorMsg:
		m.state = types.AgentStateError
		m.errText = msg.Err
	}
	return nil
}

func errEmptyInput() error {
	return errors.New("输入不能为空")
}
