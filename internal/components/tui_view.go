package components

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/admin/xagent/internal/agent/types"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("62")).Padding(0, 1)
	bannerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	roleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("36"))
	mutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

const startupBanner = `
__  __     _                    _
\ \/ /    / \   __ _  ___ _ __ | |_
 \  /    / _ \ / _  |/ _ \ '_ \| __|
 /  \   / ___ \ (_| |  __/ | | | |_
/_/\_\ /_/   \_\__, |\___|_| |_|\__|
               |___/
`

// View 渲染 TUI。布局保持朴素，方便后续扩展。
func (m *TUIModel) View() tea.View {
	width := m.width
	if width <= 0 {
		width = 80
	}

	header := titleStyle.Width(width).Render("agent-cli")
	status := m.statusLine(width)
	input := m.input.View()

	v := tea.NewView(strings.Join([]string{
		header,
		m.viewport.View(),
		input,
		status,
	}, "\n"))
	v.AltScreen = true
	return v
}

func (m *TUIModel) statusLine(width int) string {
	stateText := fmt.Sprintf("状态: %s", m.state)
	if m.errText != "" {
		stateText = errorStyle.Render(stateText + " | " + m.errText)
	} else {
		stateText = mutedStyle.Render(stateText + " | Enter 发送，Esc/Ctrl+C 退出，空输入时 q 退出")
	}
	return lipgloss.NewStyle().Width(width).Render(stateText)
}

func (m *TUIModel) syncViewport() {
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom()
}

func (m *TUIModel) renderMessages() string {
	lines := make([]string, 0, len(m.messages)*2+1)
	lines = append(lines, bannerStyle.Render(startupBanner))
	for _, message := range m.messages {
		role := renderRole(message.Role)
		content := strings.TrimSpace(message.Content)
		if content == "" {
			content = "(空消息)"
		}
		lines = append(lines, role+" "+content)
	}
	return strings.Join(lines, "\n\n")
}

func renderRole(role types.Role) string {
	switch role {
	case types.RoleUser:
		return roleStyle.Render("You")
	case types.RoleAssistant:
		return roleStyle.Render("Agent")
	case types.RoleTool:
		return roleStyle.Render("Tool")
	default:
		return mutedStyle.Render("System")
	}
}
