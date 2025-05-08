package ui

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "ghershon/internal/ui/styles"
	//"github.com/charmbracelet/bubbles/viewport"
)

type DashboardModel struct{}

func (m DashboardModel) Init() tea.Cmd {
    return tea.EnterAltScreen
}
func NewDashboardModel() DashboardModel{
	return DashboardModel{}
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, nil
}

func (m DashboardModel) View() string {
    tabs := lipgloss.JoinHorizontal(lipgloss.Top,
        styles.ActiveTab.Render("Dashboard"),
        styles.TabStyle.Render("Snippets"),
        styles.TabStyle.Render("Projects"),
        styles.TabStyle.Render("Secrets"),
        styles.TabStyle.Render("Tickets"),
    )

    //logo := styles.LogoStyle.Render("Gher ░▒▓█")
	banner := `
  .--.      .-'.      .--.      .--.      .--.      .--.      .'-.      .--.
:::::.\::::::::.\::::::::.\::::::::.\::::::::.\::::::::.\::::::::.\::::::::.\
'      '--'      '.-'      '--'      '--'      '--'      '-.'      '--'      '
░░▒▒▓▓████████████████████████████████████████████████████████▓▓▒▒░░
                       ░▒▓█ G H E R S H O N   T U I █▓▒░                       
░░▒▒▓▓████████████████████████████████████████████████████████▓▓▒▒░░

                    version 0.1.0   ·   by Hernan Gomez
              "Automate the boring. Track the cool."
`
	logo := styles.LogoStyle.Render(banner)

    content := lipgloss.JoinVertical(lipgloss.Left,
        styles.BorderBox.Width(50).Height(10).Render("Panel 1: Info"),
        styles.BorderBox.Width(50).Height(10).Render("Panel 2: Actions"),
    )

    return lipgloss.JoinVertical(lipgloss.Left, logo, tabs, content)
}
