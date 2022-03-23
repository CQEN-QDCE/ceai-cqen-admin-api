package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Retire une association d'une ressource à une autre",
	Long:  `Retire une association d'une ressource du CEAI à une autre: lab-user, lab-project, lab-account`,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
