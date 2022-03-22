package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var createProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Créer un projet Openshift",
	Long:  `Créer un projet Openshift dans l'environnement du CEAI`,
	Run: func(cmd *cobra.Command, args []string) {
		var projects []models.OpenshiftProjectWithLab

		err := HandleInput(&projects)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		for _, project := range projects {

			err = client.CreateOpenshiftProject(&project)

			if err != nil {
				fmt.Printf("Erreur: %v \n", err)
				return
			}

		}

	},
}

func init() {
	createCmd.AddCommand(createProjectCmd)
}
