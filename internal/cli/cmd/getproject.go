package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getProjectCmd = &cobra.Command{
	Use:   "project [project id]",
	Short: "Obtenir un projet Openshift",
	Long:  `Obtenir les information d'un projet Openshift`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		project, err := client.GetOpenshiftProjectFromId(args[0])

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*project)
	},
}

func init() {
	getCmd.AddCommand(getProjectCmd)
}
