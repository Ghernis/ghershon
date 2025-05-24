package ui

import (
	"fmt"
//	"os"
	"strings"
//
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
    "ghershon/internal/ui/styles"
	"ghershon/internal/storage"
	"ghershon/internal/models"
	"ghershon/internal/ui/toast"
//	"ghershon/internal/projects"
)

type SecretModel struct {
	inputs models.SecretFormInputs
	envList list.Model
	submitting bool
    focusIdx     int
    errMsg       string
	mode         *Mode
	db_service  *sql_l.SnippetsService
	toast       toast.ToastModel
}
type EnvItem struct{
	title string
	description string
}

func (i EnvItem) Title() string       { return i.title }
func (i EnvItem) Description() string { return i.description }
func (i EnvItem) FilterValue() string { return i.title }

func NewSecretModel(db_service *sql_l.SnippetsService, mode *Mode ) SecretModel{
	envs := []list.Item{
		EnvItem{"Default","Default/Global"},
		EnvItem{"DEV","Desarrollo"},
		EnvItem{"TST","Testing"},
		EnvItem{"PRD","Produccion"},
	}
	delegate := list.NewDefaultDelegate()
	envList := list.New(envs,delegate,40,5)
	envList.Title = "Environment"
	sfi:= models.SecretFormInputs{
		Name : textinput.New(),
		Env : textinput.New(),
		Desc : textinput.New(),
		Value : textinput.New(),
	    SecretType : textinput.New(),
	}
	sfi.Name.Placeholder = "Name"
	sfi.Name.Focus()
	sfi.Env.Placeholder = "Environment"
	sfi.Env.Blur()
	sfi.Desc.Placeholder = "Description"
	sfi.Desc.Blur()
	sfi.Value.Placeholder = "Value"
	sfi.Value.Blur()
	sfi.SecretType.Placeholder = "Type"
	sfi.SecretType.Blur()

	secrets:=db_service.FindAllSecret()
	projects:= db_service.FindAllProjects()
	flattened:=projects[0].Flatten()
	fmt.Println(*flattened.Description)
	fmt.Println(flattened)
	
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
		inputs: sfi,
		envList:envList,
		mode: mode,
		submitting: false,
		db_service: db_service,
		toast: toast.NewToastModel(),
	}

}

func (m SecretModel) Init() tea.Cmd{
	return tea.EnterAltScreen
}
func doSubmitSecretCmd(m SecretModel) tea.Cmd{
	return func() tea.Msg{
		//fmt.Println(m.inputs.ToSecret())
		selected := m.envList.SelectedItem().(EnvItem)
		m.inputs.Environment = selected.Title()
		err := m.db_service.AddSecret(m.inputs.ToSecret())
		if err != nil{
			return SubmitFinishedMsg{Data: "error", Err:err}
		}
		m.submitting=false
		return SubmitFinishedMsg{Data: "done", Err:nil}
	}	
}
func (m SecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	inputs := m.inputs.Slice()
	switch msg := msg.(type) {
	case SubmitFinishedMsg:
		m.submitting=false
		if msg.Err != nil{
			m.toast.Show(msg.Data  ,styles.ErrorStyle)
		}else{
			m.toast.Show("✅ Secret " +m.inputs.ToSecret().Name +" added to db!",styles.SuccessStyle)
		}
		case tea.KeyMsg:
			switch  msg.String(){
				case "i":
					if *m.mode==modeNormal{
						*m.mode=modeInsert
						return m,nil
					}
				case "esc":
					*m.mode=modeNormal	
					return m,nil
				case "enter":
					if *m.mode== modeNormal && m.submitting != true{
						m.submitting=true
						return m,tea.Batch(
							//doSubmitCmd(m.formData),
							doSubmitSecretCmd(m),
						)
					}
				case "shift+tab","backtab","up":
					if m.focusIdx > 0 {
						m.focusIdx--
					}
				case "down", "tab":
					if m.focusIdx < len(inputs) {//-1
						m.focusIdx++
						//if m.focusIdx < len(inputs)-1{
						//} else if m.focusIdx == len(inputs){	
						//} else if m.focusIdx == len(inputs)+1{	
						//}
					}
				case "k":
					if *m.mode==modeNormal{			
						if m.focusIdx > 0 {
							m.focusIdx--
						}
					}
				case "j":
					if *m.mode==modeNormal{			
						if m.focusIdx < len(inputs) {//-1
							m.focusIdx++
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
		//	if msg.String() == "ctrl+c" {
		//		return m, tea.Quit
		//	}
		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.envList.SetSize(msg.Width-h, msg.Height-v)
	 

		
	}

    cmds := make([]tea.Cmd, len(inputs)+1)
	if *m.mode!=modeNormal{
		for i := range inputs {
			inputs[i], cmds[i] = inputs[i].Update(msg)
		}
	}
	if m.focusIdx == len(inputs){
		m.envList, cmds[len(inputs)] = m.envList.Update(msg)
	}
	m.inputs.FromSlice(inputs)
    return m, tea.Batch(cmds...)
	//var cmd tea.Cmd
	//m.list, cmd = m.list.Update(msg)
	//return m, cmd
}

func (m SecretModel) View() string {
	var b strings.Builder
    b.WriteString("╔═ Create New Secret ═══════════════════════╗\n")
    for i := range m.inputs.Slice() {
        b.WriteString(fmt.Sprintf("║ %-25s %s ║\n", m.inputs.Slice()[i].Placeholder+":", m.inputs.Slice()[i].View()))
    }
    b.WriteString("╠════════════════════════════════════════════╣\n")
		b.WriteString(m.envList.View() + "\n")
    if m.errMsg != "" {
        b.WriteString(fmt.Sprintf("║ Error: %-33s ║\n", m.errMsg))
    }
	//b.WriteString(fmt.Sprintf("║ [%v] Bootstrap project files             ║\n", checkbox(m.bootstrap)))
	//b.WriteString(fmt.Sprintf("║ [%v] Create Azure DevOps ticket          ║\n", checkbox(m.createTicket)))

	b.WriteString("╠════════════════════════════════════════════╣\n")

    content := lipgloss.JoinVertical(lipgloss.Left,
        styles.BorderBox.Width(100).Height(10).Render(b.String()),
        styles.BorderBox.Width(60).Height(5).Render(m.toast.View()),
    )

    return lipgloss.JoinVertical(lipgloss.Left,   content)
    //return lipgloss.JoinVertical(lipgloss.Left, tabs, docStyle.Render(m.list.View()),footer)
}

