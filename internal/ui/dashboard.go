package ui

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
    "ghershon/internal/ui/styles"
    "ghershon/internal/models"
    "ghershon/internal/storage"
    "ghershon/internal/ui/toast"
	//"github.com/charmbracelet/bubbles/viewport"
	"strings"
	"fmt"
)


type ProjectFormModel struct {
    inputs       models.ProjectFormInputs
    focusIdx     int
    bootstrap    bool
    createTicket bool
    submitting   bool
    errMsg       string
	mode         *Mode
	db_service  *sql_l.DatabaseService
	toast       toast.ToastModel
}

func (m ProjectFormModel) Init() tea.Cmd {
    return tea.EnterAltScreen
}
func NewProjectFormModel(db_service *sql_l.DatabaseService,mode *Mode) ProjectFormModel{
	pfi:= models.ProjectFormInputs{
		Title : textinput.New(),
		Description : textinput.New(),
		Tags : textinput.New(),
		Problem_Statement : textinput.New(),
		Architecture : textinput.New(),
        Evidence : textinput.New(),
		Ticket_ID : textinput.New(),
		Expected_Finish_Date : textinput.New(),
		Completed_At : textinput.New(),
        Time_Before_Automation : textinput.New(),
		Time_After_Automation : textinput.New(),
	}

	pfi.Title.Placeholder = "Title"
	pfi.Title.Focus()
	pfi.Description.Placeholder = "Description"
	pfi.Description.Blur()
	pfi.Tags.Placeholder = "Tags"
	pfi.Tags.Blur()
	pfi.Problem_Statement.Placeholder = "Problem Statement"
	pfi.Problem_Statement.Blur()
	pfi.Architecture.Placeholder = "Architecture"
	pfi.Architecture.Blur()
	pfi.Evidence.Placeholder = "Evidence"
	pfi.Evidence.Blur()
	pfi.Ticket_ID.Placeholder = "Ticket ID"
	pfi.Ticket_ID.Blur()
	pfi.Expected_Finish_Date.Placeholder = "Expected Finish Date"
	pfi.Expected_Finish_Date.Blur()
	pfi.Completed_At.Placeholder = "Complete at"
	pfi.Completed_At.Blur()
	pfi.Time_Before_Automation.Placeholder = "Time Before Automation"
	pfi.Time_Before_Automation.Blur()
	pfi.Time_After_Automation.Placeholder = "Time After Automation"
	pfi.Time_After_Automation.Blur()
	
    return ProjectFormModel{
			inputs: pfi,
			mode: mode,
			submitting: false,
			db_service: db_service,
			toast: toast.NewToastModel(),
		}
	}

type SubmitFinishedMsg struct{
	Data string
	Err error
}


func doSubmitCmd(m ProjectFormModel) tea.Cmd{
	return func() tea.Msg{
		err := m.db_service.AddProject(m.inputs.ToProject())
		if err != nil{
			return SubmitFinishedMsg{Data: "error", Err:err}
		}
		return SubmitFinishedMsg{Data: "done", Err:nil}
	}	
}

func (m ProjectFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	inputs := m.inputs.Slice()
    switch msg := msg.(type) {
	case SubmitFinishedMsg:
		m.submitting=false
		if msg.Err != nil{
			m.toast.Show(msg.Data  ,styles.ErrorStyle)
		}else{
			m.toast.Show("✅ Project " +m.inputs.ToProject().Title +" added to db!",styles.SuccessStyle)
		}
    case tea.KeyMsg:
        switch msg.String() {
			case " ":
				if m.focusIdx == len(inputs) {
					m.bootstrap = !m.bootstrap
				} else if m.focusIdx == len(inputs)+1 {
					m.createTicket = !m.createTicket
				}
			case "enter":
				if *m.mode== modeNormal && m.submitting != true{
					m.submitting=true
					return m,tea.Batch(
						//doSubmitCmd(m.formData),
						doSubmitCmd(m),
					)
				}
			case "i":
				if *m.mode==modeNormal{
					*m.mode=modeInsert
					return m,nil
				}
			case "esc":
				*m.mode=modeNormal
				return m, nil

			case "k":
				if *m.mode==modeNormal{			
					if m.focusIdx > 0 {
						m.focusIdx--
					}
				}
			case "j":
				if *m.mode==modeNormal{			
					if m.focusIdx < len(inputs)+1 {//-1
						m.focusIdx++
						if m.focusIdx < len(inputs)-1{
						} else if m.focusIdx == len(inputs){	
						} else if m.focusIdx == len(inputs)+1{	

						}
					}
				}
			case "shift+tab","backtab","up":
				if m.focusIdx > 0 {
					m.focusIdx--
				}
			case "down", "tab":
				if m.focusIdx < len(inputs)+1 {//-1
					m.focusIdx++
					if m.focusIdx < len(inputs)-1{
					} else if m.focusIdx == len(inputs){	
					} else if m.focusIdx == len(inputs)+1{	

					}
				}
        }

        for i := range inputs {
            if i == m.focusIdx {
                inputs[i].Focus()
            } else {
                inputs[i].Blur()
            }
        }
    }

    cmds := make([]tea.Cmd, len(inputs))
	if *m.mode!=modeNormal{
		for i := range inputs {
			inputs[i], cmds[i] = inputs[i].Update(msg)
		}
	}
	m.inputs.FromSlice(inputs)
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
    for i := range m.inputs.Slice() {
        b.WriteString(fmt.Sprintf("║ %-25s %s ║\n", m.inputs.Slice()[i].Placeholder+":", m.inputs.Slice()[i].View()))
    }
    b.WriteString("╠════════════════════════════════════════════╣\n")
    if m.errMsg != "" {
        b.WriteString(fmt.Sprintf("║ Error: %-33s ║\n", m.errMsg))
    }
	//b.WriteString(fmt.Sprintf("║ [%v] Bootstrap project files             ║\n", checkbox(m.bootstrap)))
	//b.WriteString(fmt.Sprintf("║ [%v] Create Azure DevOps ticket          ║\n", checkbox(m.createTicket)))

b.WriteString("╠════════════════════════════════════════════╣\n")
b.WriteString(fmt.Sprintf("║ %-40s ║\n",
    highlightedCheckbox("Bootstrap project files", m.bootstrap, m.focusIdx == len(m.inputs.Slice())),
))
b.WriteString(fmt.Sprintf("║ %-40s ║\n",
    highlightedCheckbox("Create Azure DevOps ticket", m.createTicket, m.focusIdx == len(m.inputs.Slice())+1),
))

    b.WriteString("║ [Enter] Submit    [Esc] Cancel             ║\n")
    b.WriteString("╚════════════════════════════════════════════╝\n")
	//toast := toast.NewToastModel()


    content := lipgloss.JoinVertical(lipgloss.Left,
        styles.BorderBox.Width(60).Height(10).Render(b.String()),
        styles.BorderBox.Width(60).Height(5).Render("Panel 2: Actions"),
        styles.BorderBox.Width(60).Height(5).Render(m.toast.View()),
		
		
    )

    return lipgloss.JoinVertical(lipgloss.Left, logo, tabs, content)
}
