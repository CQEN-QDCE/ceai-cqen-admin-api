package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Obtenir une ressource",
	Long:  `Obtenir une ressource du CEAI: user[s], lab[s], project[s], account[s]`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
