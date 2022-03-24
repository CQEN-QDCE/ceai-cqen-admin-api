package cmd

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Supprimer une ressource",
	Long:  `Supprimer une ressource du CEAI: user, lab, project, account`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
