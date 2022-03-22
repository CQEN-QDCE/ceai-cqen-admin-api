package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/spf13/cobra"
)

var updateLaboratoryCmd = &cobra.Command{
	Use:     "laboratory [laboratory id]",
	Aliases: []string{"lab"},
	Short:   "Mets à jour les informations d'un laboratoire",
	Long:    `Mets à jour les informations d'un laboratoire du CEAI`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var laboratoryUpdate models.LaboratoryUpdate

		err := GetUpdateFlagsValues(&laboratoryUpdate, cmd)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		err = client.UpdateLaboratory(args[0], &laboratoryUpdate)

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}
	},
}

func init() {

	var laboratoryUpdate models.LaboratoryUpdate

	GenerateUpdateFlags(&laboratoryUpdate, updateLaboratoryCmd)

	updateCmd.AddCommand(updateLaboratoryCmd)
}
