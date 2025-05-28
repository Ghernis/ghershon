package cli

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"ghershon/internal/models"
	"strings"
)

//var (
//	secretName string
//	project int64
//	env string
//)

func PrintProjects(projects []models.Project) {
    // Header
    fmt.Printf("%-5s  %-20s  %-40s\n", "ID", "Name", "Description")
    fmt.Println(strings.Repeat("-", 80))

    // Rows
    for _, v := range projects {
		p := v.Flatten()
        fmt.Printf("%-5d  %-20s  %-40s\n",
            p.ID, p.Title, *p.Description)
    }
}
var getProjectCmd = &cobra.Command{
	Use: "get-project",
	Short: "Get a project info",
	Run: func(cmd *cobra.Command, args []string){
		if db_service.CheckInit()!= nil{
			fmt.Fprintln(os.Stderr,"Database service not initialized")
			os.Exit(1)
		}

		projects := db_service.DatabaseSrv.FindAllProjects()
		PrintProjects(projects)
	},
}

func init(){
//	getSecretCmd.Flags().StringVarP(&secretName,"name","n","","Secret name(required)")
//	getSecretCmd.Flags().Int64VarP(&project,"project","p",1,"Project name")
//	getSecretCmd.Flags().StringVarP(&env,"env","e","Default","Environment")
//
	getSecretCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(getProjectCmd)
}
