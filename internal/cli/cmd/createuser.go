package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:   "user [email nom prénom organisation role]",
	Short: "Créer un usager",
	Long:  `Créer un usager du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		var users []models.User

		err := HandleInput(&users)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		for _, user := range users {

			err = client.CreateUser(&user)

			if err != nil {
				fmt.Printf("Erreur: %v \n", err)
				return
			}

		}

	},
}

func init() {
	createCmd.AddCommand(createUserCmd)
}
