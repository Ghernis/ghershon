package utils

import(
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/go-yaml/yaml"
	"ghershon/internal/storage"
)
type Config struct{
	Bootstrap struct{
		Dir_path string `yaml:"path"`
	}	
}

func loadConfigs() Config{
	path := "./config.yml"
	data, err := os.ReadFile(path)
	if err != nil{
		fmt.Errorf("reading config file %w",err)
	}

	fmt.Println("loadConfig yml")
	conf := Config{}
	if err := yaml.Unmarshal(data,&conf); err!= nil{
		log.Fatalf("error: %v",err)
	}
	return conf

}
func DoSomething(snippetsSrv *sql_l.SnippetsService){
	log.Println("Something in db with utils")
	//snippetsSrv.GetData()

}

func Load() Config{
	fmt.Println("utils loads")
	envFile,_ := godotenv.Read(".env")

	token := envFile["TOKEN"]
	fmt.Println(token)
	return loadConfigs()
}
