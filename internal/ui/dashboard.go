package ui

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
    "ghershon/internal/ui/styles"
	//"github.com/charmbracelet/bubbles/viewport"
	"strings"
	"fmt"
)


//type DashboardModel struct{}
type ProjectFormModel struct {
    inputs       []textinput.Model
    focusIdx     int
    bootstrap    bool
    createTicket bool
    submitting   bool
    errMsg       string
	mode         *Mode
}

func (m ProjectFormModel) Init() tea.Cmd {
    return tea.EnterAltScreen
}
func NewProjectFormModel(mode *Mode) ProjectFormModel{
	fields := []string{
        "Title", "Ticket ID", "Description", "Problem Statement", "Architecture",
        "Evidence", "Expected Finish Date", "Completed At",
        "Time Before Automation", "Time After Automation", "Tags",
    }
    inputs := make([]textinput.Model, len(fields))
    for i, name := range fields {
        ti := textinput.New()
        ti.Placeholder = name
        ti.Focus() // focus first
        if i != 0 {
            ti.Blur()
        }
        inputs[i] = ti
    }
    return ProjectFormModel{
			inputs: inputs,
			mode: mode,
		}
	}

func (m ProjectFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case " ":
            if m.focusIdx == len(m.inputs) {
                m.bootstrap = !m.bootstrap
            } else if m.focusIdx == len(m.inputs)+1 {
                m.createTicket = !m.createTicket
            }
		case "enter":
			
			*m.mode=modeInsert
            if m.focusIdx == len(m.inputs)-1 {
                m.submitting = true
				m.focusIdx++
                //return m, submitToDB(m) // your function to persist to db
                return m , nil// your function to persist to db
            }
		case "esc":
			*m.mode=modeNormal
			return m, nil

        //case "up", "k":
        case "shift+tab","backtab","up":
            if m.focusIdx > 0 {
                m.focusIdx--
            }
        //case "down", "j":
        case "down", "tab":
            if m.focusIdx < len(m.inputs)+1 {//-1
				m.focusIdx++
				if m.focusIdx < len(m.inputs)-1{
				} else if m.focusIdx == len(m.inputs){	
				} else if m.focusIdx == len(m.inputs)+1{	

				}
            }
        }

        for i := range m.inputs {
            if i == m.focusIdx {
                m.inputs[i].Focus()
            } else {
                m.inputs[i].Blur()
            }
        }
    }

    cmds := make([]tea.Cmd, len(m.inputs))
    for i := range m.inputs {
        m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
    }

    return m, tea.Batch(cmds...)
}
func highlightedCheckbox(label string, checked, focused bool) string {
    box := "[ ]"
    if checked {
        box = "[x]"
    }
    if focused {
        return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(box + " " + label)
    }
    return box + " " + label
}

func checkbox(val bool) string {
    if val {
        return "x"
    }
    return " "
}

func (m ProjectFormModel) View() string {
    tabs := lipgloss.JoinHorizontal(lipgloss.Top,
        styles.ActiveTab.Render("Dashboard"),
        styles.TabStyle.Render("Snippets"),
        styles.TabStyle.Render("Bootstrap"),
        styles.TabStyle.Render("Secrets"),
        styles.TabStyle.Render("Tickets"),
    )

    //logo := styles.LogoStyle.Render("Gher ░▒▓█")
	banner := `
  .--.      .-'.      .--.      .--.      .--.      .--.      .'-.  
:::::.\::::::::.\::::::::.\::::::::.\::::::::.\::::::::.\::::::::.\:
'      '--'      '.-'      '--'      '--'      '--'      '-.'      '
░░▒▒▓▓████████████████████████████████████████████████████████▓▓▒▒░░
                  ░▒▓█ G H E R S H O N   T U I █▓▒░                       
░░▒▒▓▓████████████████████████████████████████████████████████▓▓▒▒░░

                   version 0.1.0   ·   by Hernan Gomez
	 			       "Automate the boring. Je"
`
	logo := styles.LogoStyle.Render(banner)
	var b strings.Builder
    b.WriteString("╔═ Create New Project ═══════════════════════╗\n")
    for i := range m.inputs {
        b.WriteString(fmt.Sprintf("║ %-25s %s ║\n", m.inputs[i].Placeholder+":", m.inputs[i].View()))
    }
    b.WriteString("╠════════════════════════════════════════════╣\n")
    if m.errMsg != "" {
        b.WriteString(fmt.Sprintf("║ Error: %-33s ║\n", m.errMsg))
    }
	//b.WriteString(fmt.Sprintf("║ [%v] Bootstrap project files             ║\n", checkbox(m.bootstrap)))
	//b.WriteString(fmt.Sprintf("║ [%v] Create Azure DevOps ticket          ║\n", checkbox(m.createTicket)))

b.WriteString("╠════════════════════════════════════════════╣\n")
b.WriteString(fmt.Sprintf("║ %-40s ║\n",
    highlightedCheckbox("Bootstrap project files", m.bootstrap, m.focusIdx == len(m.inputs)),
))
b.WriteString(fmt.Sprintf("║ %-40s ║\n",
    highlightedCheckbox("Create Azure DevOps ticket", m.createTicket, m.focusIdx == len(m.inputs)+1),
))

    b.WriteString("║ [Enter] Submit    [Esc] Cancel             ║\n")
    b.WriteString("╚════════════════════════════════════════════╝\n")


    content := lipgloss.JoinVertical(lipgloss.Left,
        styles.BorderBox.Width(60).Height(10).Render(b.String()),
        styles.BorderBox.Width(60).Height(10).Render("Panel 2: Actions"),
    )

    return lipgloss.JoinVertical(lipgloss.Left, logo, tabs, content)
}
