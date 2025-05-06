package styles

import "github.com/charmbracelet/lipgloss"

var (
    TabStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#7aa2f7")).
        Padding(0, 2).
        Bold(true)

    ActiveTab = TabStyle.Copy().
        Background(lipgloss.Color("#1a1b26"))

    BorderBox = lipgloss.NewStyle().
        Border(lipgloss.NormalBorder()).
        Padding(1).
        Margin(1).
        BorderForeground(lipgloss.Color("#7dcfff"))

    LogoStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#bb9af7")).
        Bold(true).
        MarginBottom(1)

	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	NoStyle      = lipgloss.NewStyle() 
)
