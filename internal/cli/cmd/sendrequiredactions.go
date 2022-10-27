package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var sendRequiredActionsCmd = &cobra.Command{
	Use:     "required-actions [user]",
	Aliases: []string{"ra"},
	Short:   "Envois un courriel à un usager pour compléter ses actions requises. ",
	Long:    `Envois un courriel à un usager du CEAI pour compléter ses actions requises, comme la création d'un mot de passe ou d'un jeton OTP`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.SendRequiredActionEmail(args[0])

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}
	},
}

func init() {

	sendCmd.AddCommand(sendRequiredActionsCmd)
}
