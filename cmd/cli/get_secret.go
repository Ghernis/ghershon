package cli

import (
	"fmt"
	"os"
	"strings"

	"ghershon/internal/security"
	"ghershon/internal/models"
	"github.com/spf13/cobra"
)

var (
	secretName string
	project int64
	env string
	listAll bool
)

func listSecrets(){
	PrintSecrets(db_service.DatabaseSrv.FindAllSecret())
}

func PrintSecrets(secrets []models.Secret) {
    // Header
    fmt.Printf("%-5s  %-20s  %-20s  %-20s\n", "ID", "Name", "Env","Project")
    fmt.Println(strings.Repeat("-", 80))

    // Rows
    for _, p := range secrets {
		//p := v.Flatten()
        fmt.Printf("%-5d  %-20s  %-20s  %-20v\n",
            p.Id, p.Name,p.Environment, p.Project_id)
    }
}

var getSecretCmd = &cobra.Command{
	Use: "get-secret",
	Short: "Get a secret value",
	Run: func(cmd *cobra.Command, args []string){
		if db_service.CheckInit()!= nil{
			fmt.Fprintln(os.Stderr,"Database service not initialized")
			os.Exit(1)
		}
		if listAll{
			listSecrets()
			os.Exit(0)
		}

		
		secret := db_service.DatabaseSrv.FindSecretFiltered(secretName,project,env)
		if (len(secret)>0){
			val,err := encrypt.DecryptText(secret[0].Encoded_value,db_service.Key_secret)
			if err != nil{
				fmt.Fprintln(os.Stderr,"Error decrypting value")
			}
			fmt.Println(val)
		}
	},
}


func init(){
	getSecretCmd.Flags().StringVarP(&secretName,"name","n","","Secret name(required)")
	getSecretCmd.Flags().BoolVarP(&listAll,"list","l",false,"List all secrets")
	getSecretCmd.Flags().Int64VarP(&project,"project","p",1,"Project name")
	getSecretCmd.Flags().StringVarP(&env,"env","e","Default","Environment")

	//getSecretCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(getSecretCmd)
}
