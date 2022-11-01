package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Désauthentification",
	Long:  `Détruit la session en cours au serveur d'API`,
	Run: func(cmd *cobra.Command, args []string) {
		Logout()
	},
}

func init() {
	rootCmd.AddCommand(LogoutCmd)
}

func Logout() {
	err := client.DeleteSession()

	if err != nil {
		fmt.Printf("Erreur: %v \n", err)
		return
	}

	fmt.Printf("Session terminée. \n")
}
