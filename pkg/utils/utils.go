package utils

import(
	"fmt"
	"log"
	"github.com/joho/godotenv"
	//"github.com/go-yaml/yaml"
	"ghershon/internal/storage"
)

func loadConfigs(){
	fmt.Println("loadConfig yml")

}
func DoSomething(snippetsSrv *sql_l.SnippetsService){
	log.Println("Something in db with utils")
	snippetsSrv.GetData()

}

func Load(){
	fmt.Println("utils loads")
	envFile,_ := godotenv.Read(".env")

	token := envFile["TOKEN"]
	fmt.Println(token)
	loadConfigs()
}
