package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var deleteUserCmd = &cobra.Command{
	Use:   "user [email]",
	Short: "Supprime un usager",
	Long:  `Supprimer un usager du CEAI`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.DeleteUser(args[0])

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		fmt.Printf("L'usager %v supprimé", args[0])
	},
}

func init() {

	deleteCmd.AddCommand(deleteUserCmd)
}
