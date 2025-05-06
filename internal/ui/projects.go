package ui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
    "ghershon/internal/ui/styles"
	"ghershon/internal/ui/toast"
	"ghershon/internal/projects"
	"ghershon/pkg/utils"
)

type ProjectTypeItem struct{
	title string
	description string
}

func (i ProjectTypeItem) Title() string       { return i.title }
func (i ProjectTypeItem) Description() string { return i.description }
func (i ProjectTypeItem) FilterValue() string { return i.title }

type BootstrapModel struct {
	inputs     []textinput.Model
	//projectName     textinput.Model
	//projectName     textinput.Model
	//projectName     textinput.Model
	projectList list.Model
	showalist bool
	focusIndex int
	projectType string
	quitting    bool
	toast toast.ToastModel
}

func NewBootstrapModel() BootstrapModel {

	items := []list.Item{
		ProjectTypeItem{"Python Script","Simple CLI or Pipe"},
		ProjectTypeItem{"Django + Helm","Django micro with Helm Chart"},
		ProjectTypeItem{"Go CLI","Go Command Line"},
		ProjectTypeItem{"React","React bootstrap"},
	}
	delegate := list.NewDefaultDelegate()
	projectList := list.New(items,delegate,40,5)
	projectList.Title = "Project Type"

	var inputs []textinput.Model

	// Project Name
	projectName := textinput.New()
	projectName.Placeholder = "Project name"
	projectName.Prompt = "Name: "
	projectName.Focus()
	projectName.PromptStyle = styles.FocusedStyle
	projectName.TextStyle = styles.FocusedStyle

	// Location
	location := textinput.New()
	location.Placeholder = "~/dev/"
	location.Prompt = "Path: "

	// Registry
	registry := textinput.New()
	registry.Placeholder = "ghcr.io/user/"
	registry.Prompt = "Registry: "

	inputs = []textinput.Model{
		projectName, 
		location,
		registry,
	}

	return BootstrapModel{
		projectList: projectList,
		inputs:      inputs,
		focusIndex:  0,
		projectType: "Python Script", // Fake dropdown, static for now
		toast: toast.NewToastModel(),
	}
}

func (m BootstrapModel) Init() tea.Cmd {
	return textinput.Blink
}
func saveDebug(name, path, registry, kind string) {
	f, _ := os.Create("debug.log")
	defer f.Close()
	fmt.Fprintf(f, "Name: %s\nPath: %s\nRegistry: %s\nType: %s\n", name, path, registry, kind)
}

func (m BootstrapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.focusIndex == len(m.inputs)-1 {
				// Last input, submit
				name := m.inputs[0].Value()
				path := m.inputs[1].Value()
				registry := m.inputs[2].Value()
				selected := m.projectList.SelectedItem().(ProjectTypeItem)

				saveDebug(
					name,
					path,
					registry,
					selected.title,
				)
				if(selected.title == "Python Script"){
					config:=utils.Load()
					bootstrap.Python_boot(config.Bootstrap.Dir_path,name)
				}
				if(selected.title == "Django + Helm"){
					config:=utils.Load()
					bootstrap.Django_boot(config.Bootstrap.Dir_path,name)
				}
				//m.toast.Show("✅ Snippet saved!", styles.SuccessStyle)
				return m, tea.Quit //quit?
			}
			// Move focus to next
			m.focusIndex++
		case "up":
			m.focusIndex--
			if m.focusIndex <0{
				m.focusIndex = len(m.inputs)-1
			}	
			m= m.updateFocus()
		case "tab","down":
			m.focusIndex++
			if m.focusIndex >= len(m.inputs){
				m.focusIndex=0
			}
			m= m.updateFocus()
			
		}

	case tea.WindowSizeMsg:
		// Optional: handle resizing
	}
	m.toast,_ = m.toast.Update(msg)

	// Forward to correct component
	if m.focusIndex == 0 {
		m.projectList, cmd = m.projectList.Update(msg)
	} else {
		m.inputs[m.focusIndex-1], cmd = m.inputs[m.focusIndex-1].Update(msg)
	}

	return m, cmd
}
func (m BootstrapModel) updateFocus() BootstrapModel {
	for i := range m.inputs {
		if i == m.focusIndex-1 {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = styles.FocusedStyle
			m.inputs[i].TextStyle = styles.FocusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = styles.NoStyle
			m.inputs[i].TextStyle = styles.NoStyle
		}
	}
	return m
}

//func (m *BootstrapModel) updateFocus() (tea.Model, tea.Cmd) {
//	cmds := make([]tea.Cmd, len(m.inputs))
//	for i := 0; i < len(m.inputs); i++ {
//		if i == m.focusIndex {
//			m.inputs[i].Focus()
//			m.inputs[i].PromptStyle = styles.FocusedStyle
//			m.inputs[i].TextStyle = styles.FocusedStyle
//		} else {
//			m.inputs[i].Blur()
//			m.inputs[i].PromptStyle = styles.NoStyle
//			m.inputs[i].TextStyle = styles.NoStyle
//		}
//	}
//	return *m, tea.Batch(cmds...)
//}

func (m BootstrapModel) View() string {
	//m=m.initialModel()
	//m.Init()
    tabs := lipgloss.JoinHorizontal(lipgloss.Top,
        styles.TabStyle.Render("Dashboard"),
        styles.TabStyle.Render("Snippets"),
        styles.ActiveTab.Render("Projects"),
    )
	if m.quitting {
		return " Exiting...\n"
	}

	var b strings.Builder

	fmt.Fprintln(&b, tabs)
	fmt.Fprintln(&b, "\n┌──────────────────── New Project ───────────────────┐")

	if m.focusIndex == 0 {
		b.WriteString(m.projectList.View() + "\n")
	} else {
		selected := m.projectList.SelectedItem().(ProjectTypeItem)
		fmt.Fprintf(&b, "│ Project Type: %s\n", selected.title)
	}

	for _, input := range m.inputs {
		fmt.Fprintf(&b, "│ %s\n", input.View())
	}
	if toastMsg :=m.toast.View(); toastMsg != ""{
		fmt.Fprintf(&b, "| %s\n",toastMsg)
	}

	fmt.Fprintln(&b, "│")
	fmt.Fprintln(&b, "│ [Tab] Next • [Enter] Confirm • [Esc] Cancel")
	fmt.Fprintln(&b, "└────────────────────────────────────────────────────┘")
	return b.String()
}


//func main() {
//	if err := tea.NewProgram(initialModel()).Start(); err != nil {
//		fmt.Println("Error:", err)
//		os.Exit(1)
//	}
//}

