package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var addLaboratoryUserCmd = &cobra.Command{
	Use:     "laboratory-user [laboratory id] [user1 user2 userx...]",
	Aliases: []string{"lab-user"},
	Short:   "Associe un usager à un laboratoire",
	Long:    `Associe un usager à un laboratoire du CEAI`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.AddLaboratoryUsers(args[0], args[1:])

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}
	},
}

func init() {

	addCmd.AddCommand(addLaboratoryUserCmd)
}
