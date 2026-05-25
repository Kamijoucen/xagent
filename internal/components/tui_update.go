package components

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/admin/xagent/internal/agent/types"
)

// Update 处理 Bubble Tea 消息。框架回调属于 TUI Component，不归类为 Logic。
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd := m.applyMessage(msg); cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = max(1, msg.Width)
		m.viewport.Width = msg.Width
		m.viewport.Height = max(3, msg.Height-6)
		m.syncViewport()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "q":
			if strings.TrimSpace(m.input.Value()) == "" {
				return m, tea.Quit
			}
		case "enter":
			return m, m.handleSubmit()
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *TUIModel) handleSubmit() tea.Cmd {
	value := strings.TrimSpace(m.input.Value())
	if value == "" {
		m.ReportError(errEmptyInput())
		return nil
	}

	m.messages = append(m.messages, types.Message{Role: types.RoleUser, Content: value, CreatedAt: time.Now()})
	m.input.SetValue("")
	m.state = types.AgentStateThinking
	m.errText = ""
	m.syncViewport()

	if handler := m.currentSubmitHandler(); handler != nil {
		return handler(value)
	}
	return nil
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}
