package sql_l

//bd: snippets

import (
	"html/template"
	"time"
	"bytes"
	"strings"
	//"database/sql"
	"fmt"
	"os"
	"log"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewDB(driver, dsn string) (*sqlx.DB, error) {
	// Create if not exist
	if _, err := os.Stat(dsn); os.IsNotExist(err) {
		log.Println("Database does not exist. Creating...")
		file, err := os.Create(dsn)
		if err != nil {
			log.Fatalf("Failed to create database file: %v", err)
		}
		file.Close()
	}
	// Connect
    db, err := sqlx.Connect(driver, dsn)
    if err != nil {
        return nil, err
    }
	// Initialize
	err = initializeSchema(db)
	if err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	log.Println("Database initialized successfully.")
    return db, nil
}
func initializeSchema(db *sqlx.DB) error{
	path := "internal/storage/tables.sql"
	schema, err := os.ReadFile(path)
	if err != nil{
		fmt.Errorf("Reading db file %w",err)
	}
	statements := strings.Split(string(schema),";")
	for	 _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_,err = db.Exec(stmt)
		if err != nil{
			fmt.Errorf("Error in statement: %w",err)
		}
	}	

	//err = seedProjectsTable(db)
	if err != nil{
		fmt.Errorf("Seeding projects %w",err)
	}
	return err
}

func seedProjectsTable(db *sqlx.DB) error {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM projects")
	if err != nil {
		return err
	}
	if count == 0 {
		_, err := db.Exec(`
			INSERT INTO projects (title, description, created_at)
			VALUES (?, ?, CURRENT_TIMESTAMP)
		`, "Example Project", "This is a starter project")
		return err
	}
	return nil
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
	Id_ticket             string    `db:"id_ticket"`
	Description           string    `db:"description"`
	ProblemStatement      string    `db:"problem_statement"`
	Architecture          string    `db:"architecture"`
	Evidence              string    `db:"evidence"`
	ExpectedFinishDate    string    `db:"expected_finish_date"`
	CompletedAt           *string   `db:"completed_at"`
	TimeBeforeAutomation  int       `db:"time_before_automation"`
	TimeAfterAutomation   int       `db:"time_after_automation"`
	Tags                  string    `db:"tags"`
	//CreatedAt             time.Time `db:"created_at"`
	CreatedAt             string    `db:"created_at"`
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
type Secret struct {
	Id int64 `db:"id"`
	Name string `db:"name"`
	Description string `db:"description"`
	Secret_type string `db:"secret_type"`
	Encoded_value string `db:"encoded_value"`
	//Created_at time.Time `db:"created_at"`
	Created_at string `db:"created_at"`
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
func (s *SnippetsService) FindSecret(search_id int64) []Secret{
	//fmt.Println(search_string)
	getData := fmt.Sprintf(`
		select * from secrets where  id == %v;
	`,search_id)
	var data []Secret
	fmt.Println(getData)
	err := s.db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}
func (s *SnippetsService) FindAllSecret()[]Secret{
	//fmt.Println(search_string)
	getData := fmt.Sprintf(`
		select * from secrets;
	`)
	var data []Secret
	fmt.Println(getData)
	err := s.db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}

func (s *SnippetsService) AddSecret(sec Secret) error{
	//fmt.Println(search_string)
	query := fmt.Sprintf(`
        INSERT INTO secrets (name, description, secret_type, encoded_value)
		VALUES (:name,:description, :secret_type,  :encoded_value)
		`)

	_,err := s.db.NamedExec(query,sec)
	
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return err
}

func (s *SnippetsService) AddProject(project Project) error{
	query := fmt.Sprintf(`
        INSERT INTO projects (title,id_ticket, description, problem_statement, architecture,evidence,expected_finish_date,completed_at,time_before_automation,time_after_automation,tags)
		VALUES (:title,id_ticket, :description, :problem_statement, :architecture,:evidence,:expected_finish_date,:completed_at,:time_before_automation,:time_after_automation,:tags)
		`)

	_,err := s.db.NamedExec(query,project)
	
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return err
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
