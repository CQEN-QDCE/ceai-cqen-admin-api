package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset [password|otp|all] [user]",
	Short: "Réinitialise le crédentiel d'un usager",
	Long:  `Réinitialise le crédentiel d'un usager`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.ResetUserCredential(args[1], args[0])

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		fmt.Printf("Demande de réinitialisation envoyée à %v", args[1])
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
