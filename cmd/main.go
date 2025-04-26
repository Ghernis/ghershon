package main

import(
	//"fmt"
	"log"
	"ghershon/internal/ui"
	"ghershon/internal/storage"
	"ghershon/pkg/utils"

	"github.com/charmbracelet/bubbletea"

	"github.com/jmoiron/sqlx"
	//_ "modernc.org/sqlite"
)

type App struct{
	db *sqlx.DB
	SnippetsSrv *sql_l.SnippetsService
}

func newApp() *App{
	db:= sql_l.MustNewDB("sqlite","./ghershon.db")

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &App{
		db: db,
		SnippetsSrv: sql_l.NewSnippetsService(db),
	}
}
func main(){
	//ui.Load()
	//sql.Load()
	utils.Load()

	app := newApp()
	defer app.db.Close()
	app.SnippetsSrv.GetData()
	utils.DoSomething(app.SnippetsSrv)
	p := tea.NewProgram(ui.RootModel{})
	if err := p.Start(); err != nil{
		panic(err)
	}


}
