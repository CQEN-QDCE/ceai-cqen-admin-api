package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getLaboratoryCmd = &cobra.Command{
	Use:     "laboratory [lab id]",
	Aliases: []string{"lab"},
	Short:   "Obtenir un laboratoire",
	Long:    `Obtenir les information d'un laboratoire du CEAI`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lab, err := client.GetLaboratoryFromId(args[0])

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*lab)
	},
}

func init() {
	getCmd.AddCommand(getLaboratoryCmd)
}
