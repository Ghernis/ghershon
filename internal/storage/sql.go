package sql_l

//bd: snippets

import (
	"html/template"
	"time"
	//"strconv"
	"bytes"
	_ "embed"
	"strings"
	//"database/sql"
    "ghershon/internal/models"
	"ghershon/internal/security"
	"fmt"
	"os"
	"log"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed tables.sql
var schemaSQL string

type DatabaseService struct {
	Db *sqlx.DB
}


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
		log.Fatalf("Failed to initialize schema in NewDB: %v", err)
	}
    return db, nil
}
func initializeSchema(db *sqlx.DB) error{
	//path := "internal/storage/tables.sql"
	statements := strings.Split(schemaSQL, ";")
	for _,stmt := range statements{
		stmt = strings.TrimSpace(stmt)
		if stmt==""{
			continue
		}
		_,err := db.Exec(stmt)
		if err != nil{
			return fmt.Errorf("Error in SQL statement: %w",err)
		}
	}
	
//	schema, err := os.ReadFile(path)
//	if err != nil{
//		fmt.Errorf("Reading db file %w",err)
//	}
//	statements := strings.Split(string(schema),";")
//	for	 _, stmt := range statements {
//		stmt = strings.TrimSpace(stmt)
//		if stmt == "" {
//			continue
//		}
//		_,err = db.Exec(stmt)
//		if err != nil{
//			fmt.Errorf("Error in statement: %w",err)
//		}
//	}	

	if err := seedProjectsTable(db); err !=nil{
		return fmt.Errorf("Seeding projects %w",err)
	}
	return nil
}

func seedProjectsTable(db *sqlx.DB) error {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM projects")
	if err != nil {
		return err
	}
	if count == 0 {
		desc := "This is the global project config"
		_, err := db.Exec(`
			INSERT INTO projects (title, description )
			VALUES (?, ?)
		`, "Global Project", &desc)
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


func NewDatabaseService(db *sqlx.DB) *DatabaseService {
	return &DatabaseService{Db: db}
}

//func (p ProjectDTO) toDB() Project{//viendo
//	return Project{
//		Title:                p.Title,
//		Id_ticket:            derefString(p.Id_ticket),
//		Description:          derefString(p.Description),
//		ProblemStatement:     derefString(p.ProblemStatement),
//		Architecture:         derefString(p.Architecture),
//		Evidence:             derefString(p.Evidence),
//		ExpectedFinishDate:   derefString(p.ExpectedFinishDate),
//		CompletedAt:          derefString(p.CompletedAt),
//		TimeBeforeAutomation: derefInt(p.TimeBeforeAutomation),
//		TimeAfterAutomation:  derefInt(p.TimeAfterAutomation),
//		Tags:                 derefString(p.Tags),
//		CreatedAt:            p.CreatedAt,
//		
//	}.Flatten()
//}


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

func (s *DatabaseService) FindData2(search_string string) []Snippet{
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
	err := s.Db.Select(&data,query.String())
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}

func (s *DatabaseService) FindData(search_string string) []Snippet{
	
	getData := fmt.Sprintf(`
	select * from snippets where description like '%%%s%%';
	`,search_string)
	var data []Snippet
	fmt.Println(getData)
	err := s.Db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}
func (s *DatabaseService) FindSecret(search_id int64) []models.Secret{
	//fmt.Println(search_string)
	getData := fmt.Sprintf(`
		select * from secrets where  id == %v;
	`,search_id)
	var data []models.Secret
	fmt.Println(getData)
	err := s.Db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}
func (s *DatabaseService) FindSecretFiltered(name string,project_id int64,env string) []models.Secret{
	//fmt.Println(search_string)
	getData := fmt.Sprintf(`
		select * from secrets where  (name == '%v' and project_id == %v and environment == '%v');
	`,name,project_id,env)
	var data []models.Secret
	//fmt.Println(getData)
	err := s.Db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}
func (s *DatabaseService) FindAllSecret()[]models.Secret{
	//fmt.Println(search_string)
	getData := fmt.Sprintf(`
		select * from secrets;
	`)
	var data []models.Secret
	fmt.Println(getData)
	err := s.Db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}

func (s *DatabaseService) AddSecret(sec models.Secret,key_secret []byte) error{
	val,err := encrypt.EncryptText(sec.Encoded_value,key_secret)
	if err != nil{
		fmt.Fprintln(os.Stderr,"Error decrypting value")
	}
	sec.Encoded_value = val
	
	query := fmt.Sprintf(`
        INSERT INTO secrets (name,project_id,environment, description, secret_type, encoded_value,is_encrypted,created_at)
		VALUES (:name,:project_id,:environment,:description, :secret_type, :encoded_value,:is_encrypted,:created_at)
		`)

	_,err = s.Db.NamedExec(query,sec)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return err
}

func (s *DatabaseService) FindAllProjects()[]models.Project{
	//fmt.Println(search_string)
	getData := fmt.Sprintf(`
		select * from projects;
	`)
	var data []models.Project
	//fmt.Println(getData)
	err := s.Db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data
}

func (s *DatabaseService) AddProject(project models.Project) error{
	query := fmt.Sprintf(`
        INSERT INTO projects (title,id_ticket, description, problem_statement, architecture,evidence,expected_finish_date,completed_at,time_before_automation,time_after_automation,tags)
		VALUES (:title,:id_ticket, :description, :problem_statement, :architecture,:evidence,:expected_finish_date,:completed_at,:time_before_automation,:time_after_automation,:tags)
		`)

	_,err := s.Db.NamedExec(query,project)
	
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

func (s *DatabaseService) GetData() ([]Snippet,error){
	getData := `
	select * from snippets;   
	`
	var data []Snippet
	err := s.Db.Select(&data,getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	return data,err
}
