package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"ghershon/internal/storage"
)

type Screen int

const (
    Dashboard Screen = iota
    Snippets
    Bootstrap
    Secret
)

type RootModel struct {
    current Screen
    dash    DashboardModel
	snippets SnippetModel
	bootstrap BootstrapModel
	secret SecretModel
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
        case "4":
            m.current = Secret
			return m,m.secret.Init()
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
		case Secret:
			var c tea.Cmd
			var updated tea.Model
			updated, c=m.secret.Update(msg)
			m.secret = updated.(SecretModel)
			cmd=c
	}
    return m, cmd
}
func NewRootModel(db_service *sql_l.SnippetsService) RootModel{
	return RootModel{
		current: Dashboard,
		dash: NewDashboardModel(),
		snippets: NewSnippetModel(),
		bootstrap: NewBootstrapModel(),
		secret: NewSecretModel(db_service),
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
    case Secret:
        return m.secret.View()

    default:
        return "Unknown screen"
    }
}
