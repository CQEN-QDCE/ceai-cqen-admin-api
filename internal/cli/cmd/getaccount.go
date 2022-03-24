package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getAccountCmd = &cobra.Command{
	Use:   "account [account id]",
	Short: "Obtenir un compte AWS",
	Long:  `Obtenir les information d'compte AWS`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account, err := client.GetAwsAccount(args[0])

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*account)
	},
}

func init() {
	getCmd.AddCommand(getAccountCmd)
}
