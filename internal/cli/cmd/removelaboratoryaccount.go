package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var removeLaboratoryAccountCmd = &cobra.Command{
	Use:     "laboratory-account [laboratory id] [account id]",
	Aliases: []string{"lab-account"},
	Short:   "Retire l'association d'un compte AWS à un laboratoire",
	Long:    `Retire l'association d'un compte AWS à un laboratoire du CEAI`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.DetachAwsAccountFromLaboratory(args[0], args[1])

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}
	},
}

func init() {

	removeCmd.AddCommand(removeLaboratoryAccountCmd)
}
