package models

import (
	"strconv"
	"github.com/charmbracelet/bubbles/textinput"
)

type SecretFormInputs struct{
        Name textinput.Model
		Env textinput.Model
		Desc textinput.Model
		Value textinput.Model
		SecretType textinput.Model
		Environment string
}
func (i SecretFormInputs) Slice() []textinput.Model {
	return []textinput.Model{
		i.Name,
		i.Env,
		i.Desc,
		i.Value,
		i.SecretType,
	}
}

type Secret struct {
	Id int64 `db:"id"`
    Project_id int64 `db:"project_id"`
    Environment string `db:"environment"`
	Name string `db:"name"`
	Description string `db:"description"`
	Secret_type string `db:"secret_type"`
	Encoded_value string `db:"encoded_value"`
    Is_encrypted bool `db:"is_encrypted"`
	//Created_at time.Time `db:"created_at"`
	Created_at string `db:"created_at"`
}
func (p SecretFormInputs) ToSecret() Secret {
    return Secret{
        Name:                p.Name.Value(),
		Project_id: 0,
        Description:          p.Desc.Value(),
		Environment:          p.Environment,
        Secret_type:         p.SecretType.Value(),
        Encoded_value:             p.Value.Value(),
		Is_encrypted: true,
        Created_at:            "", // or generate default timestamp here
    }
}
func (i *SecretFormInputs) FromSlice(inputs []textinput.Model){
	i.Name = inputs[0]
	i.Env = inputs[1]
	i.Desc = inputs[2]
	i.Value = inputs[3]
	i.SecretType = inputs[4]
}

type ProjectFormInputs struct{
        Title textinput.Model
		Description textinput.Model
		Tags textinput.Model
		Problem_Statement textinput.Model
		Architecture textinput.Model
        Evidence textinput.Model
		Ticket_ID textinput.Model
		Expected_Finish_Date textinput.Model
		Completed_At textinput.Model
        Time_Before_Automation textinput.Model
		Time_After_Automation textinput.Model
}

func (i ProjectFormInputs) Slice() []textinput.Model {
	return []textinput.Model{
		i.Title,
		i.Description,
		i.Tags,
		i.Problem_Statement,
		i.Architecture,
        i.Evidence,
		i.Ticket_ID,
		i.Expected_Finish_Date,
		i.Completed_At,
        i.Time_Before_Automation,
		i.Time_After_Automation,
	}
}

type Project struct {
	ID                    int64     `db:"id"`
	Title                 string    `db:"title"`
	Id_ticket             *string    `db:"id_ticket"`
	Description           *string    `db:"description"`
	ProblemStatement      *string    `db:"problem_statement"`
	Architecture          *string    `db:"architecture"`
	Evidence              *string    `db:"evidence"`
	ExpectedFinishDate    *string    `db:"expected_finish_date"`
	CompletedAt           *string   `db:"completed_at"`
	TimeBeforeAutomation  *int       `db:"time_before_automation"`
	TimeAfterAutomation   *int       `db:"time_after_automation"`
	Tags                  *string    `db:"tags"`
	//CreatedAt             time.Time `db:"created_at"`
	CreatedAt             *string    `db:"created_at"`

}

func (p ProjectFormInputs) ToProject() Project {
    return Project{
        Title:                p.Title.Value(),
        Id_ticket:            toPtr(p.Ticket_ID.Value()),
        Description:          toPtr(p.Description.Value()),
        ProblemStatement:     toPtr(p.Problem_Statement.Value()),
        Architecture:         toPtr(p.Architecture.Value()),
        Evidence:             toPtr(p.Evidence.Value()),
        ExpectedFinishDate:   toPtr(p.Expected_Finish_Date.Value()),
        CompletedAt:          toPtr(p.Completed_At.Value()),
        TimeBeforeAutomation: toIntPtr(p.Time_Before_Automation.Value()),
        TimeAfterAutomation:  toIntPtr(p.Time_After_Automation.Value()),
        Tags:                 toPtr(p.Tags.Value()),
        CreatedAt:            nil, // or generate default timestamp here
    }
}
func toPtr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}
func toIntPtr(s string) *int {
    if i, err := strconv.Atoi(s); err == nil {
        return &i
    }
    return nil
}

func (i *ProjectFormInputs) FromSlice(inputs []textinput.Model){
	i.Title = inputs[0]
	i.Description = inputs[1]
	i.Tags = inputs[2]
	i.Problem_Statement = inputs[3]
	i.Architecture = inputs[4]
	i.Evidence = inputs[5]
	i.Ticket_ID = inputs[6]
	i.Expected_Finish_Date = inputs[7]
	i.Completed_At = inputs[8]
	i.Time_Before_Automation = inputs[9]
	i.Time_After_Automation = inputs[10]
}

func (p Project) Flatten() Project {
	return Project{
		ID:                   p.ID,
		Title:                p.Title,
		Id_ticket:            derefString(p.Id_ticket),
		Description:          derefString(p.Description),
		ProblemStatement:     derefString(p.ProblemStatement),
		Architecture:         derefString(p.Architecture),
		Evidence:             derefString(p.Evidence),
		ExpectedFinishDate:   derefString(p.ExpectedFinishDate),
		CompletedAt:          derefString(p.CompletedAt),
		TimeBeforeAutomation: derefInt(p.TimeBeforeAutomation),
		TimeAfterAutomation:  derefInt(p.TimeAfterAutomation),
		Tags:                 derefString(p.Tags),
		CreatedAt:            p.CreatedAt,
	}
}
func derefString(s *string) *string {
	if s == nil {
		empty := ""
		return &empty
	}
	return s
}

func derefInt(i *int) *int {
	if i == nil {
		zero := 0
		return &zero
	}
	return i
}
