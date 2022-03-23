package cmd

import "github.com/spf13/cobra"

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Associe une ressource à une autre",
	Long:  `Associe une ressource du CEAI à une autre: lab-user, lab-project, lab-account`,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
