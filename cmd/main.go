package main

import(
	//"fmt"
	"log"
	"os"
		
	"ghershon/cmd/cli"
	"ghershon/internal/ui"
	"ghershon/internal/storage"
	"ghershon/internal/security"
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
	dbPath, err  := utils.GetDataPath("ghershon.db")
	if err !=nil{
		log.Fatal("Failet to get DB path: ",err)
	}
	db:= sql_l.MustNewDB("sqlite",dbPath)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &App{
		db: db,
		SnippetsSrv: sql_l.NewSnippetsService(db),
	}
}
func main(){

	_,err:=encrypt.EnsureEncryptionKey()
	if err != nil{
		log.Fatal(err)
	}

	app := newApp()
	defer app.db.Close()
	if len(os.Args) >1{
		cli.Execute(app.SnippetsSrv)
	} else{
		//utils.DoSomething(app.SnippetsSrv)
		p := tea.NewProgram(ui.NewRootModel(app.SnippetsSrv))
		if err := p.Start(); err != nil{
			panic(err)
		}
	}
}

