package ui

import (
	//"fmt"
//	"os"
//	"strings"
//
	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
    "ghershon/internal/ui/styles"
	"ghershon/internal/storage"
//	"ghershon/internal/ui/toast"
//	"ghershon/internal/projects"
)

type SecretModel struct {
	list list.Model
}
//type item struct {
//	title, desc string
//}
//
//func (i item) Title() string       { return i.title }
//func (i item) Description() string { return i.desc }
//func (i item) FilterValue() string { return i.title }

func NewSecretModel(db_service *sql_l.SnippetsService ) SecretModel{
	secrets:=db_service.FindAllSecret()

	var items []list.Item
	for _,v := range secrets{
		items=append(items,item{
			title:v.Name,
			desc:v.Description,
		})
	}
	mylist:= list.New(items, list.NewDefaultDelegate(), 40, 40)
	mylist.Title = "Secrets"
	return SecretModel{
		list:mylist,
	}

}

func (m SecretModel) Init() tea.Cmd{
	return tea.EnterAltScreen
}
func (m SecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//case tea.KeyMsg: // TODO cambiar esto
	//	if msg.String() == "ctrl+c" {
	//		return m, tea.Quit
	//	}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SecretModel) View() string {
    tabs := lipgloss.JoinHorizontal(lipgloss.Top,
        styles.TabStyle.Render("Dashboard"),
        styles.TabStyle.Render("Snippets"),
        styles.TabStyle.Render("Bootstrap"),
        styles.ActiveTab.Render("Secrets"),
        styles.TabStyle.Render("Tickets"),
    )
	footer := "Press [← →] to change month · [enter] Chart view · [q] Quit"
    return lipgloss.JoinVertical(lipgloss.Left, tabs, docStyle.Render(m.list.View()),footer)
}

