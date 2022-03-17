package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Obtenir la liste des usagers",
	Long:  `Obtenir la liste des usagers du CEAI`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		users, err := client.GetUsers()

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*users)
	},
}

func init() {
	getCmd.AddCommand(getUsersCmd)
}
