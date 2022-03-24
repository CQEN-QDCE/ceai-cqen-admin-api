package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var updateUserCmd = &cobra.Command{
	Use:   "user [email]",
	Short: "Créer un usager",
	Long:  `Créer un usager du CEAI`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var userUpdate models.UserUpdate

		err := GetUpdateFlagsValues(&userUpdate, cmd)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		err = client.UpdateUser(args[0], &userUpdate)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}
	},
}

func init() {

	var userUpdate models.UserUpdate

	GenerateUpdateFlags(&userUpdate, updateUserCmd)

	updateCmd.AddCommand(updateUserCmd)
}
