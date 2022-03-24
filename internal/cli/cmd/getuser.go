package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getUserCmd = &cobra.Command{
	Use:   "user [username]",
	Short: "Obtenir ou des usagers",
	Long:  `Obtenir ou des usagers du CEAI`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user, err := client.GetUserFromUsername(args[0])

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*user)
	},
}

func init() {
	getCmd.AddCommand(getUserCmd)
}
