package utils

import(
	"fmt"
	"log"
	"os"
	//"github.com/joho/godotenv"
	"github.com/go-yaml/yaml"
	"ghershon/internal/storage"
	"path/filepath"
	"encoding/base64"
)
type Config struct{
	Bootstrap struct{
		Dir_path string `yaml:"path"`
	}	
}

func GetConfigPath(filename string) (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "Error getting user config dir", err
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
		return "Error getting user config dir", err
	}
	appDir := filepath.Join(dataDir,"ghershon") 
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, filename), nil
}

func LoadKey(filename string) ([]byte, error) {
	dataDir, err := os.UserConfigDir() 
	if err != nil {
		return nil, err
	}
	keyDir := filepath.Join(dataDir,"ghershon",filename) 
    encoded, err := os.ReadFile(keyDir)
    if err != nil {
        return nil, err
    }
    key, err := base64.StdEncoding.DecodeString(string(encoded))
    if err != nil {
        return nil, err
    }
    if len(key) != 32 {
        return nil, fmt.Errorf("key must be 32 bytes")
    }
    return key, nil
}

func loadConfigs() Config{
	//path := "./config.yml"
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
func DoSomething(snippetsSrv *sql_l.DatabaseService){
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
