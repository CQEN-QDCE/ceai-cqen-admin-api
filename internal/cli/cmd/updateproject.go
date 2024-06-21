package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var updateProjectCmd = &cobra.Command{
	Use:   "project [project id]",
	Short: "Mets à jour les informations d'un projet Openshift",
	Long:  `Mets à jour les informations d'un projet Openshift du CEAI`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectUpdate models.OpenshiftProjectUpdate

		err := GetUpdateFlagsValues(&projectUpdate, cmd)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		err = client.UpdateOpenshiftProject(args[0], &projectUpdate)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		fmt.Printf("Project Openshift %v mis à jour", args[0])
	},
}

func init() {

	var projectUpdate models.OpenshiftProjectUpdate

	GenerateUpdateFlags(&projectUpdate, updateProjectCmd)

	updateCmd.AddCommand(updateProjectCmd)
}
