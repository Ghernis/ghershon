package main

import(
	//"fmt"
	"log"
	"os"
		
	"ghershon/cmd/cli"
	"ghershon/internal/ui"
	"ghershon/internal/appstate"
	"ghershon/internal/storage"
	"ghershon/internal/security"
	"ghershon/pkg/utils"

	"github.com/charmbracelet/bubbletea"

	"github.com/jmoiron/sqlx"
	//_ "modernc.org/sqlite"
)

type App struct{
	db *sqlx.DB
	store appstate.AppServices
}

func newApp(key []byte) *App{
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
		store: appstate.AppServices{
			DatabaseSrv: sql_l.NewDatabaseService(db),
			KeySecret: key,
		},
	}
}
func main(){

	key_name,err:=encrypt.EnsureEncryptionKey()
	if err != nil{
		log.Fatal(err)
	}
	key_encrypt,err := utils.LoadKey(key_name)
	if err != nil{
		log.Fatal(err)
	}

	app := newApp(key_encrypt)
	defer app.db.Close()

	if len(os.Args) >1{
		cli.Execute(app.store.DatabaseSrv,app.store.KeySecret)
	} else{
		//utils.DoSomething(app.SnippetsSrv)
		p := tea.NewProgram(ui.NewRootModel(app.store.DatabaseSrv))
		if err := p.Start(); err != nil{
			panic(err)
		}
	}
}

