package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var addLaboratoryProjectCmd = &cobra.Command{
	Use:     "laboratory-project [laboratory id] [project id]",
	Aliases: []string{"lab-project"},
	Short:   "Associe un projet Openshift à un laboratoire",
	Long:    `Associe un projet Openshift à un laboratoire du CEAI`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.AttachOpenshiftProjectToLaboratory(args[0], args[1])

		if err != nil {
			fmt.Printf("Erreur: %v \n", err)
			return
		}

		fmt.Printf("Project Openshift #%v associé au laboratoire %v", args[1], args[0])
	},
}

func init() {

	addCmd.AddCommand(addLaboratoryProjectCmd)
}
