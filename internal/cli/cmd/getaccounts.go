package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Obtenir la liste des comptes AWS",
	Long:  `Obtenir la liste des comptes AWS du CEAI`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		accounts, err := client.GetAwsAccounts()

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*accounts)
	},
}

func init() {
	getCmd.AddCommand(getAccountsCmd)
}
