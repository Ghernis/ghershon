package main

import(
	"fmt"
	"ghershon/internal/ui"
	"ghershon/internal/storage"
	"ghershon/pkg/utils"
	"log"

	"database/sql"
	_ "modernc.org/sqlite"
)

type App struct{
	db *sql.DB
	SnippetSrv *sql_l.SnippetService
}

func newApp() *App{
	// Connect to or create the database
	db, err := sql.Open("sqlite", "./mydata.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Connected to SQLite!")
	//snippetService := NewSnippetService(db)
	//fmt.Println("Created service")
	return &App{
		db: db,
		SnippetSrv: sql_l.NewSnippetService(db),
	}
}
func main(){
	fmt.Println("hello")
	ui.Load()
	//sql.Load()
	utils.Load()

	app := newApp()
	defer app.db.Close()
	app.SnippetSrv.GetData()
	utils.DoSomething(app.SnippetSrv)

}
