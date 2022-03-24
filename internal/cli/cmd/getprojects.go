package cmd

import (
	"fmt"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/cli/client"
	"github.com/spf13/cobra"
)

var getProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Obtenir la liste des projets Openshift",
	Long:  `Obtenir la liste des projets Openshift du CEAI`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		projects, err := client.GetOpenshiftProjects()

		if err != nil {
			fmt.Printf("Error: %v \n", err)
			return
		}

		PrintOutput(*projects)
	},
}

func init() {
	getCmd.AddCommand(getProjectsCmd)
}
