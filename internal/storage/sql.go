package sql_l

import (
	"database/sql"
	"fmt"
	"log"

	//_ "modernc.org/sqlite"
)

// Service struct that depends on the database
type SnippetService struct {
	db *sql.DB
}
// Constructor function
func NewSnippetService(db *sql.DB) *SnippetService {
	return &SnippetService{db: db}
}

//func initDb() *SnippetService{
//	// Connect to or create the database
//	db, err := sql.Open("sqlite", "./mydata.db")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	// Test the connection
//	if err := db.Ping(); err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Connected to SQLite!")
//	snippetService := NewSnippetService(db)
//	fmt.Println("Created service")
//	return snippetService
//}
func (s *SnippetService) GetData(){
	getData := `
	select * from snippets;`
	data, err := s.db.Query(getData)
	if err != nil {
		log.Fatal("Error in query: ",err)
	}
	//defer data.Close()

	for data.Next(){
		var (
			id int64
			title string
			description string
			code string
			tags string
			reference string
		)
		if err := data.Scan(&id,&title,&description,&code,&tags,&reference); err != nil{
			log.Fatal(err)
			fmt.Println(err)
		}
		log.Printf("ID: %d, Title: %s, Description: %s,%s,%s,%s\n", id, title, description,code,tags,reference)
		//fmt.Println(id,title,description,code,tags,reference)
	}

}

//func Load() {

	//ss := initDb()

	//ss.getData()

	// Create a sample table
	//createTableSQL := `
	//CREATE TABLE IF NOT EXISTS snippets (
	//	id INTEGER PRIMARY KEY AUTOINCREMENT,
	//	title TEXT NOT NULL,
	//	description TEXT,
	//	code TEXT,
	//	tags TEXT,
	//	reference TEXT)
	//	`

	//_, err = db.Exec(createTableSQL)
	//if err != nil {
	//	log.Fatalf("Error creating table: %s", err)
	//}
	//fmt.Println(res)

	//fmt.Println("Table created successfully!")

	//putData := `
	//	INSERT INTO snippets (title, description, code, tags,reference)
	//	VALUES ('titulo','description','code','tags','ref');
	//`
	//_,err = db.Exec(putData)
	//if err != nil{
	//	log.Fatal("Insert Error:",err)
	//}
	//}
//}
