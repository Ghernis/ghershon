package ui

import (
	"fmt"
	//"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
    "ghershon/internal/ui/styles"
)

type BootstrapModel struct {
	inputs     []textinput.Model
	focusIndex int
	projectType string
	quitting    bool
}

func (BootstrapModel)initialModel() BootstrapModel {
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
		inputs:      inputs,
		focusIndex:  0,
		projectType: "Python Script", // Fake dropdown, static for now
	}
}

func (m BootstrapModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m BootstrapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.focusIndex == len(m.inputs)-1 {
				// Last input, submit
				return m, tea.Quit
			}
			// Move focus to next
			m.focusIndex++
		case "up":
			m.focusIndex--
			if m.focusIndex <0{
				m.focusIndex = len(m.inputs)-1
			}	
			return m.updateFocus()
		case "tab","down":
			m.focusIndex++
			if m.focusIndex >= len(m.inputs){
				m.focusIndex=0
			}
			return m.updateFocus()
			
		}

	case tea.WindowSizeMsg:
		// Optional: handle resizing
	}

	// Handle input updates
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m *BootstrapModel) updateFocus() (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = styles.FocusedStyle
			m.inputs[i].TextStyle = styles.FocusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = styles.NoStyle
			m.inputs[i].TextStyle = styles.NoStyle
		}
	}
	return *m, tea.Batch(cmds...)
}

func (m BootstrapModel) View() string {
	m=m.initialModel()
	m.Init()
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
	fmt.Fprintln(&b, "\n┌─────────────── New Project ────────────────┐")
	fmt.Fprintf(&b, "│ Project Type: %s\n", m.projectType)
	fmt.Fprintln(&b, "│")
	for _, input := range m.inputs {
		fmt.Fprintf(&b, "│ %s\n", input.View())
	}
	fmt.Fprintln(&b, "│")
	fmt.Fprintln(&b, "│ [Enter] Next  [Esc] Cancel")
	fmt.Fprintln(&b, "└────────────────────────────────────────────┘")
	return b.String()
}


//func main() {
//	if err := tea.NewProgram(initialModel()).Start(); err != nil {
//		fmt.Println("Error:", err)
//		os.Exit(1)
//	}
//}

