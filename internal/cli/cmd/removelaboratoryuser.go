package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var removeLaboratoryUserCmd = &cobra.Command{
	Use:     "laboratory-user [laboratory id] [user1 user2 userx...]",
	Aliases: []string{"lab-user"},
	Short:   "Retire l'association d'un usager à un laboratoire",
	Long:    `Retire l'association d'un usager à un laboratoire du CEAI`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.RemoveLaboratoryUsers(args[0], args[1:])

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		fmt.Printf("%v usager(s) desassocié(s) du laboratoire %v", len(args[1:]), args[0])
	},
}

func init() {

	removeCmd.AddCommand(removeLaboratoryUserCmd)
}
