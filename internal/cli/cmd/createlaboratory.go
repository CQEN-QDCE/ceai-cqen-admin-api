package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var createLaboratoryCmd = &cobra.Command{
	Use:     "laboratory",
	Aliases: []string{"lab"},
	Short:   "Créer un laboratoire",
	Long:    `Créer un laboratoire dans l'environnement du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		var laboratories []models.Laboratory

		err := HandleInput(&laboratories)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		for _, laboratory := range laboratories {

			err = client.CreateLaboratory(&laboratory)

			if err != nil {
				fmt.Printf("Erreur: %v \n", err)
				return
			}

		}

	},
}

func init() {
	createCmd.AddCommand(createLaboratoryCmd)
}
