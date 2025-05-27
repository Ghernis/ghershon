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
var db_service *sql_l.SnippetsService

func Execute(service *sql_l.SnippetsService){
	db_service = service
	if err := rootCmd.Execute(); err != nil{
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
}
