package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	header := ""
	body := baseStyle.Render(m.table.View()) + "\n"
	footer := m.footerView()
	return header + body + footer
}

func (m model) footerView() string {
	if m.err != nil {
		return m.errorView(m.err)
	}
	return "\n"
}

func (m model) errorView(err error) string {
	return "\n" + fmt.Sprintf("发生错误了: %s，按 q 按键退出", err.Error())
}
