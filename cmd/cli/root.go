package cli

import(
		"fmt"
		"os"
		"github.com/spf13/cobra"
		"ghershon/internal/storage"
)

var rootCmd = &cobra.Command{
	Use: "ghershon",
	Short: "Headless CLI for Ghershon",
	Long: "Use Ghershon in CLI mode to fetch secrets, bootstrap projects and more.\n\tAuthor: Hernan Gomez",
}

type DB_Service struct{
	DatabaseSrv *sql_l.DatabaseService
	Key_secret []byte
}

func (ds DB_Service) CheckInit() error{
	if ds.DatabaseSrv == nil {
        return fmt.Errorf("database service is not initialized")
    }
    if ds.Key_secret == nil {
        return fmt.Errorf("encryption key is not initialized")
    }
    return nil
}

var db_service DB_Service

func Execute(service *sql_l.DatabaseService,key []byte){
	db_service = DB_Service{
		DatabaseSrv: service,
		Key_secret: key,
	}
	if err := rootCmd.Execute(); err != nil{
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
}
