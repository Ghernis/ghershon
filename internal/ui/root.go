package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ghershon/internal/storage"
	"fmt"
	"strings"
)

type Screen int
type Mode int

const (
    ProjectForm Screen = iota
    Snippets
    Bootstrap
    Secret
)
const (
	modeInsert Mode = iota
	modeNormal 
)

type RootModel struct {
    current Screen
    dash    ProjectFormModel
	snippets SnippetModel
	bootstrap BootstrapModel
	secret SecretModel
	mode *Mode
}

func (m RootModel) Init() tea.Cmd {

	//optional of the first screen to preload
	cmd := m.bootstrap.Init()
    return tea.Batch(tea.EnterAltScreen,cmd)
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
		if *m.mode == modeNormal{
			switch msg.String(){
				case "q", "ctrl+c":
					return m,tea.Quit
			}
		}
        switch msg.String() {
        case "1":
            m.current = ProjectForm
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
		case ProjectForm:
			var c tea.Cmd
			var updated tea.Model
			updated, c =m.dash.Update(msg)
			m.dash = updated.(ProjectFormModel)
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
	mode := modeNormal
	return RootModel{
		current: ProjectForm,
		dash: NewProjectFormModel(&mode),
		snippets: NewSnippetModel(&mode),
		bootstrap: NewBootstrapModel(),
		secret: NewSecretModel(db_service,&mode),
		mode: &mode,
	}	
}
func screenName(s Screen) string {
    switch s {
    case ProjectForm:
        return "Project Form"
    case Snippets:
        return "Snippets"
    case Bootstrap:
        return "Bootstrap"
    case Secret:
        return "Secrets"
    default:
        return "Unknown"
    }
}

func modeName(m Mode) string {
    if m == modeNormal {
        return "NORMAL"
    }
    return "INSERT"
}

func (m RootModel) View() string {
	var content string
    switch m.current {
    case ProjectForm:
        content = m.dash.View()
    case Snippets:
        content = m.snippets.View()
    case Bootstrap:
        content = m.bootstrap.View()
    case Secret:
        content = m.secret.View()


    default:
        content = "Unknown screen"
    }
	header := fmt.Sprintf("╔═ Tab: %s ══════════════════════════════════╗", screenName(m.current))
    footer := fmt.Sprintf("║ Mode: %s · [q] quit · [tab] move ║", modeName(*m.mode))
    separator := "╠" + strings.Repeat("═", 44) + "╣"
	return lipgloss.JoinVertical(lipgloss.Left,
        header,
        separator,
        content,
        separator,
        footer,
    )
}
