package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var WhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Informations sur session en cours",
	Long:  `Décrit la session en cours au serveur d'API`,
	Run: func(cmd *cobra.Command, args []string) {
		Whoami()
	},
}

func init() {
	rootCmd.AddCommand(WhoamiCmd)
}

func Whoami() {
	authUser, err := client.Whoami()

	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}

	session, err := client.GetSession()

	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}

	authUser.Server = &session.ServerUrl

	PrintOutput(*authUser)
}
