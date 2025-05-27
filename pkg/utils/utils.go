package utils

import(
	"fmt"
	"log"
	"os"
	//"github.com/joho/godotenv"
	"github.com/go-yaml/yaml"
	"ghershon/internal/storage"
	"path/filepath"
)
type Config struct{
	Bootstrap struct{
		Dir_path string `yaml:"path"`
	}	
}

func GetConfigPath(filename string) (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(cfgDir, "ghershon")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, filename), nil
}

func GetDataPath(filename string) (string, error) {
	dataDir, err := os.UserConfigDir() 
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dataDir,"ghershon") 
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, filename), nil
}

func loadConfigs() Config{
	path := "./config.yml"
	path,err := GetConfigPath("config.yml")
	if err != nil{
		fmt.Errorf("Error finding config file %w",err)
	}
	data, err := os.ReadFile(path)
	if err != nil{
		fmt.Errorf("reading config file %w",err)
	}

	//fmt.Println("loadConfig yml")
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
	//fmt.Println("utils loads")
	//envFile,_ := godotenv.Read(".env")

	//token := envFile["TOKEN"]
	//fmt.Println(token)
	return loadConfigs()
}
