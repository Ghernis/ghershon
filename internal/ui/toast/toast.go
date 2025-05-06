package toast

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ghershon/internal/ui/styles"
)	

type ToastModel struct{
	visible bool
	message string
	style lipgloss.Style
	duration time.Duration
	timer time.Time
}

func NewToastModel() ToastModel{
	return ToastModel{
		visible: false,
		duration: 3* time.Second,
		style: styles.InfoStyle,
	}
}

func (m *ToastModel) Show(msg string, style lipgloss.Style){
	m.message = msg
	m.style = style
	m.visible = true
	m.timer = time.Now()
}

func (m *ToastModel) Update(msg tea.Msg) (ToastModel,tea.Cmd){
	if m.visible && time.Since(m.timer) < m.duration {
		m.visible = false
	}
	return *m,nil
}

func (m ToastModel) View() string{
		if !m.visible {
			return ""
		}
		return m.style.Render(m.message)
}
