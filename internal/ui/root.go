package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
    Dashboard Screen = iota
    Snippets
    Bootstrap
)

type RootModel struct {
    current Screen
    dash    DashboardModel
	snippets SnippetModel
}

func (m RootModel) Init() tea.Cmd {
    return tea.EnterAltScreen
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
		case "q", "ctrl+c":
			return m,tea.Quit
        case "1":
            m.current = Dashboard
        case "2":
            m.current = Snippets
        case "3":
            m.current = Bootstrap
        }
    }
    return m, nil
}

func (m RootModel) View() string {
    switch m.current {
    case Dashboard:
        return m.dash.View()
    case Snippets:
        return m.snippets.View()

    default:
        return "Unknown screen"
    }
}
