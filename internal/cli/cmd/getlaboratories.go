package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getLaboratoriesCmd = &cobra.Command{
	Use:     "laboratories",
	Aliases: []string{"labs"},
	Short:   "Obtenir la liste des laboratoires",
	Long:    `Obtenir la liste des laboratoires du CEAI`,
	Args:    cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		labs, err := client.GetLaboratories()

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*labs)
	},
}

func init() {
	getCmd.AddCommand(getLaboratoriesCmd)
}
