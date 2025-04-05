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
func DoSomething(snippetSrv *sql_l.SnippetService){
	log.Println("Something in db with utils")
	snippetSrv.GetData()

}

func Load(){
	fmt.Println("utils loads")
	envFile,_ := godotenv.Read(".env")

	token := envFile["TOKEN"]
	fmt.Println(token)
	loadConfigs()
}
