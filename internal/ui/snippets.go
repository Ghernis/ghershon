package ui

import (
	"fmt"

	//"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/viewport"
    "ghershon/internal/ui/styles"
	//"ghershon/internal/storage"
)

// Styling TODO: to ui/styles
var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69"))
	docStyle 	   = lipgloss.NewStyle().Margin(5, 2)
	headerStyle    = lipgloss.NewStyle().Bold(true).Underline(true)
	sectionTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("33"))
	itemStyle      = lipgloss.NewStyle().PaddingLeft(2)
	highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	boxStyle       = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("60"))//125 pink maybe
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type SnippetModel struct {
	list list.Model
	viewport viewport.Model
	content string
	ready bool
	mode *Mode
}

func (m SnippetModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}
//func loadData(snippetsSrv *sql_l.SnippetsService) []sql_l.Secret{
	//secrets:=app.SnippetsSrv.FindAllSecret()
//	return secrets
//}

func NewSnippetModel(mode *Mode) SnippetModel{
	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
		item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		item{title: "Linux", desc: "Pretty much the best OS"},
		item{title: "Business school", desc: "Just kidding"},
		item{title: "Pottery", desc: "Wet clay is a great feeling"},
		item{title: "Shampoo", desc: "Nothing like clean hair"},
		item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		item{title: "Stickers", desc: "The thicker the vinyl the better"},
		item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}
	


	mylist:= list.New(items, list.NewDefaultDelegate(), 50, 50)
	mylist.Title = "My Fave Things"
	return SnippetModel{
		list:mylist,
		mode: mode,
	}
}


func (m SnippetModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String(){
			case "q":
				if *m.mode == modeInsert{
					return m, nil
				}
				return m,nil
			case "esc":
				//fmt.Println("esc en snipets")
				*m.mode=modeNormal
				return m, nil
			case "enter":
				*m.mode=modeInsert
				return m, nil
	}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}


func (m SnippetModel) View() string {
    tabs := lipgloss.JoinHorizontal(lipgloss.Top,
        styles.ActiveTab.Render("Dashboard"),
        styles.TabStyle.Render("Snippets"),
        styles.TabStyle.Render("Bootstrap"),
        styles.TabStyle.Render("Secrets"),
        styles.TabStyle.Render("Tickets"),
    )

	title := titleStyle.Render(" Personal Finance Dashboard (2024)")
	header := fmt.Sprintf("Hola: %s", "hernan")
	footer := "Press [← →] to change month · [enter] Chart view · [q] Quit"
	//m=m.InitialModel()

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		tabs,
		docStyle.Render(m.list.View()),
		headerStyle.Render(header),
		highlightStyle.Render(footer),
	)
}
