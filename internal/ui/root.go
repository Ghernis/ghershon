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
	bootstrap BootstrapModel
}

func (m RootModel) Init() tea.Cmd {

	//optional of the first screen to preload
	cmd := m.bootstrap.Init()
    return tea.Batch(tea.EnterAltScreen,cmd)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
		case "q", "ctrl+c":
			return m,tea.Quit
        case "1":
            m.current = Dashboard
			return m,m.dash.Init()
        case "2":
            m.current = Snippets
			return m,m.snippets.Init()
        case "3":
            m.current = Bootstrap
			return m,m.bootstrap.Init()
        }
    }
	var cmd tea.Cmd
	switch m.current{
		case Dashboard:
			var c tea.Cmd
			var updated tea.Model
			updated, c =m.dash.Update(msg)
			m.dash = updated.(DashboardModel)
			cmd=c
		case Snippets:
			var c tea.Cmd
			var updated tea.Model
			updated, c=m.snippets.Update(msg)
			m.snippets = updated.(SnippetModel)
			cmd=c
		case Bootstrap:
			var c tea.Cmd
			var updated tea.Model
			updated, c=m.bootstrap.Update(msg)
			m.bootstrap = updated.(BootstrapModel)
			cmd=c
	}
    return m, cmd
}
func NewRootModel() RootModel{
	return RootModel{
		current: Dashboard,
		dash: NewDashboardModel(),
		snippets: NewSnippetModel(),
		bootstrap: NewBootstrapModel(),
	}	
}

func (m RootModel) View() string {
    switch m.current {
    case Dashboard:
        return m.dash.View()
    case Snippets:
        return m.snippets.View()
    case Bootstrap:
        return m.bootstrap.View()

    default:
        return "Unknown screen"
    }
}
