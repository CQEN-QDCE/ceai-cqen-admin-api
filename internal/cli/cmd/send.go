package cmd

import "github.com/spf13/cobra"

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Envoyer un message a un usager",
	Long:  `Envoyer un message a un usager du CEAI: required-actions`,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
