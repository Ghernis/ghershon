package cli

import (
	"fmt"
	"os"

	//"ghershon/internal/storage"
	"github.com/spf13/cobra"
)

var (
	secretName string
	project int64
	env string
)

var getSecretCmd = &cobra.Command{
	Use: "get-secret",
	Short: "Get a secret value",
	Run: func(cmd *cobra.Command, args []string){
		if db_service.CheckInit()!= nil{
			fmt.Fprintln(os.Stderr,"Database service not initialized")
			os.Exit(1)
		}

		
		secret := db_service.DatabaseSrv.FindSecretFiltered(secretName,project,env)
		if (len(secret)>0){
			fmt.Println(secret[0].Encoded_value)
		}
	},
}

func init(){
	getSecretCmd.Flags().StringVarP(&secretName,"name","n","","Secret name(required)")
	getSecretCmd.Flags().Int64VarP(&project,"project","p",1,"Project name")
	getSecretCmd.Flags().StringVarP(&env,"env","e","Default","Environment")

	getSecretCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(getSecretCmd)
}
