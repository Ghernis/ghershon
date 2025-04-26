package sql_l

//bd: snippets

import (
	"html/template"
	"time"
	"bytes"
	//"database/sql"
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewDB(driver, dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Connect(driver, dsn)
    if err != nil {
        return nil, err
    }
    return db, nil
}

func MustNewDB(driver, dsn string) *sqlx.DB {
    db, err := NewDB(driver, dsn)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

type SnippetsService struct {
	db *sqlx.DB
}

func NewSnippetsService(db *sqlx.DB) *SnippetsService {
	return &SnippetsService{db: db}
}

type Project struct {
	ID                    int64     `db:"id"`
	Title                 string    `db:"title"`
	Description           string    `db:"description"`
	ProblemStatement      string    `db:"problem_statement"`
	Architecture          string    `db:"architecture"`
	Evidence              string    `db:"evidence"`
	ExpectedFinishDate    string    `db:"expected_finish_date"`
	CompletedAt           *string   `db:"completed_at"`
	TimeBeforeAutomation  int       `db:"time_before_automation"`
	TimeAfterAutomation   int       `db:"time_after_automation"`
	Tags                  string    `db:"tags"`
	CreatedAt             time.Time `db:"created_at"`
}

type ProjectTask struct {
	ID         int64     `db:"id"`
	ProjectID  int64     `db:"project_id"`
	Content    string    `db:"content"`
	IsDone     bool      `db:"is_done"`
	CreatedAt  time.Time `db:"created_at"`
	DueDate    *string   `db:"due_date"`
}

type Snippet struct {
	ID             int64     `db:"id"`
	Title          string    `db:"title"`               // Short label
	Description    string    `db:"description"`         // What it does or why it's useful
	Language       string    `db:"language"`            // "yaml", "go", "bash", etc.
	Tags           string    `db:"tags"`                // CSV or JSON encoded slice
	Content        string    `db:"content"`             // Actual snippet text
	SourceFile     string    `db:"source_file"`         // Optional: file path
	StartLine      int       `db:"start_line"`          // Optional: for reference only
	EndLine        int       `db:"end_line"`            // Optional: for reference only
	Documentation  string    `db:"documentation_url"`   // Optional: official links
	ProjectUsedIn  string    `db:"project_used_in"`     // For traceability
	//CreatedAt      time.Time `db:"created_at"`
	CreatedAt      string `db:"created_at"`
	ProjectID  int64     `db:"project_id"`
}

func (s *SnippetsService) FindData2(search_string string) []Snippet{
	type search_query struct{
		Search string
	}
	
	getData := `
	select * from snippets where description like '%{{ .Search }}%';
	`
	t := template.Must(template.New("test").Parse(getData))
	query := new(bytes.Buffer)
	t.Execute(query,search_query{search_string})
	var data []Snippet
	fmt.Println(query.String())
	err := s.db.Select(&data,query.String())
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}

func (s *SnippetsService) FindData(search_string string) []Snippet{
	
	getData := fmt.Sprintf(`
	select * from snippets where description like '%%%s%%';
	`,search_string)
	var data []Snippet
	fmt.Println(getData)
	err := s.db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}

//type Sums struct{
//	Mes string `db:"mes"`
//	Year string `db:"year"`
//	Tot float64 `db:"Tot"`
//	Card string `db:"card"`
//
//}	
//func (s *SnippetsService) GetMonthlyExpenses() []Sums{
//	
//	getData := fmt.Sprintf(`
//	select year,Mes,SUM(pesos) as Tot,card from snippets group by mes,year,card;
//	`)
//	var data []Sums
//	err := s.db.Select(&data,getData)
//	if err != nil {
//		log.Fatal("Error in query: ",err)
//	}
//	fmt.Println(data)
//	return data
//}

func (s *SnippetsService) GetData() ([]Snippet,error){
	getData := `
	select * from snippets;   
	`
	var data []Snippet
	err := s.db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data,err
}
